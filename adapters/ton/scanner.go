// Scanner provides functionality for scanning the TON blockchain.
// It includes logic for fetching blocks, transactions, and accounts.
// Some network handling logic was partially taken from
// https://github.com/xssnick/ton-payment-network/blob/master/tonpayments/chain/block-scan.go (thanks to the author).
package ton

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	addressutils "github.com/xssnick/tonutils-go/address"
	tlbutils "github.com/xssnick/tonutils-go/tlb"
	tonutils "github.com/xssnick/tonutils-go/ton"
	"golang.org/x/sync/errgroup"

	"github.com/kriuchkov/tonbeacon/pkg/retrier"
)

const (
	// defaultWaitNodeTimeout is the default timeout for waiting for a node to be available.
	defaultWaitNodeTimeout = 20 * time.Second
)

type accFetchTask struct {
	master   *tonutils.BlockIDExt
	shard    *tonutils.BlockIDExt
	tx       *tlbutils.Transaction
	addr     *addressutils.Address
	callback func()
}

type OptionsScanner struct {
	NumWorkers int
}

func (o *OptionsScanner) SetDefaults() {
	if o.NumWorkers == 0 {
		o.NumWorkers = 10
	}
}

type Scanner struct {
	retrier *retrier.Retrier

	api        APIClientWrapped
	lastBlock  uint32
	numWorkers int
	taskPool   chan accFetchTask
}

func NewScanner(api APIClientWrapped, opt *OptionsScanner) *Scanner {
	opt.SetDefaults()

	return &Scanner{
		api:        api,
		taskPool:   make(chan accFetchTask, 100),
		numWorkers: opt.NumWorkers,
		retrier:    retrier.NewRetrier(),
	}
}

func (v *Scanner) RunAsync(ctx context.Context, ch chan<- any) error {
	go v.accFetcherWorker(ch, v.numWorkers)

	master, err := v.api.GetMasterchainInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "get masterchain info")
	}

	log.Debug().Uint32("seqno", master.SeqNo).Msg("starting scanner")

	masters := []*tonutils.BlockIDExt{master}
	go func() {
		outOfSync := false
		for {
			start := time.Now()

			var transactionsNum, shardBlocksNum uint64
			wg := sync.WaitGroup{}
			wg.Add(len(masters))
			for _, m := range masters {
				go func(m *tonutils.BlockIDExt) {
					defer wg.Done()

					txNum, bNum := v.fetchBlock(context.Background(), m)
					atomic.AddUint64(&transactionsNum, txNum)
					atomic.AddUint64(&shardBlocksNum, bNum)
				}(m)
			}

			wg.Wait()

			took := time.Since(start)
			log.Debug().Uint32("seqno", masters[len(masters)-1].SeqNo).Dur("took", took).Msg("scanned master")

			lastProcessed := masters[len(masters)-1]
			blocksNum := len(masters)
			masters = masters[:0]

			for {
				lastMaster, err := v.api.WaitForBlock(lastProcessed.SeqNo + 1).GetMasterchainInfo(ctx)
				if err != nil {
					log.Debug().Err(err).Uint32("seqno", lastProcessed.SeqNo+1).Msg("failed to get last block")
					continue
				}

				if lastMaster.SeqNo <= lastProcessed.SeqNo {
					continue
				}

				diff := lastMaster.SeqNo - lastProcessed.SeqNo
				if diff > 60 {
					rd := took.Round(time.Millisecond)
					if shardBlocksNum > 0 {
						rd /= time.Duration(shardBlocksNum)
					}

					log.Warn().Uint32("lag_master_blocks", diff).
						Int("processed_master_blocks", blocksNum).
						Uint64("processed_shard_blocks", shardBlocksNum).
						Uint64("processed_transactions", transactionsNum).
						Dur("took_ms_per_block", rd).
						Msg("chain scanner is out of sync")

					outOfSync = true
				} else if diff <= 1 && outOfSync {
					log.Info().Msg("chain scanner is synchronized")
					outOfSync = false
				}

				log.Debug().Uint32("lag_master_blocks", diff).
					Uint64("processed_transactions", transactionsNum).Msg("scanner delay")

				if diff > 100 {
					diff = 100
				}

				for i := lastProcessed.SeqNo + 1; i <= lastProcessed.SeqNo+diff; i++ {
					for {
						nextMaster, err := v.api.WaitForBlock(i).LookupBlock(ctx, lastProcessed.Workchain, lastProcessed.Shard, i)
						if err != nil {
							log.Debug().Err(err).Uint32("seqno", i).Msg("failed to get next block")
						}

						masters = append(masters, nextMaster)
						break
					}
				}

				v.lastBlock = masters[len(masters)-1].SeqNo
				break
			}
		}
	}()
	return nil
}

func (v *Scanner) accFetcherWorker(ch chan<- any, threads int) {
	for range threads {
		go func() {
			for {
				task := <-v.taskPool

				func() {
					defer task.callback()

					var acc *tlbutils.Account
					{
						ctx := context.Background()
						for range 20 {
							var err error
							ctx, err = v.api.Client().StickyContextNextNode(ctx)
							if err != nil {
								log.Debug().Err(err).Str("addr", task.addr.String()).Msg("failed to pick next node")
								break
							}

							qCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
							acc, err = v.api.WaitForBlock(task.master.SeqNo).GetAccount(qCtx, task.master, task.addr)
							cancel()
							if err != nil {
								log.Debug().Err(err).Str("addr", task.addr.String()).Msg("failed to get account")
								time.Sleep(100 * time.Millisecond)
								continue
							}
							break
						}
					}

					if acc == nil || !acc.IsActive || acc.State.Status != tlbutils.AccountStatusActive {
						return
					}

					ch <- task.tx
				}()
			}
		}()
	}
}

func (v *Scanner) getNotSeenShards(
	ctx context.Context,
	api APIClientWrapped,
	shard *tonutils.BlockIDExt,
	prevShards []*tonutils.BlockIDExt,
) (ret []*tonutils.BlockIDExt, err error) {
	if shard.Workchain != 0 {
		return nil, nil
	}

	if slices.ContainsFunc(prevShards, shard.Equals) {
		return nil, nil
	}

	b, err := api.GetBlockData(ctx, shard)
	if err != nil {
		return nil, errors.Wrap(err, "get block data")
	}

	parents, err := b.BlockInfo.GetParentBlocks()
	if err != nil {
		return nil, errors.Wrap(err, "get parent blocks")
	}

	for _, parent := range parents {
		ext, err := v.getNotSeenShards(ctx, api, parent, prevShards)
		if err != nil {
			return nil, errors.Wrap(err, "get not seen shards")
		}
		ret = append(ret, ext...)
	}
	return append(ret, shard), nil
}

func (v *Scanner) fetchBlock(ctx context.Context, master *tonutils.BlockIDExt) (transactionsNum, shardBlocksNum uint64) {
	log.Debug().Uint32("seqno", master.SeqNo).Msg("scanning master")

	tm := time.Now()
	for {
		select {
		case <-ctx.Done():
			log.Warn().Uint32("master", master.SeqNo).Msg("ctx done")
			return transactionsNum, shardBlocksNum
		default:
		}

		prevMaster, err := v.api.WaitForBlock(master.SeqNo-1).LookupBlock(ctx, master.Workchain, master.Shard, master.SeqNo-1)
		if err != nil {
			log.Debug().Err(err).Uint32("seqno", master.SeqNo-1).Msg("failed to get prev master block")
			continue
		}

		prevShards, err := v.api.GetBlockShardsInfo(ctx, prevMaster)
		if err != nil {
			log.Debug().Err(err).Uint32("master", master.SeqNo).Msg("failed to get shards on block")
			continue
		}

		currentShards, err := v.api.GetBlockShardsInfo(ctx, master)
		if err != nil {
			log.Debug().Err(err).Uint32("master", master.SeqNo).Msg("failed to get shards on block")
			continue
		}

		log.Debug().Uint32("seqno", master.SeqNo).Dur("took", time.Since(tm)).Msg("shards fetched")

		var newShards []*tonutils.BlockIDExt
		for _, shard := range currentShards {
			for {
				select {
				case <-ctx.Done():
					return transactionsNum, shardBlocksNum
				default:
				}

				var notSeen []*tonutils.BlockIDExt
				notSeen, err = v.getNotSeenShards(ctx, v.api, shard, prevShards)
				if err != nil {
					log.Debug().Err(err).Uint32("master", master.SeqNo).Msg("get not seen shards on block")
					continue
				}

				newShards = append(newShards, notSeen...)
				break
			}
		}

		atomic.AddUint64(&shardBlocksNum, uint64(len(newShards)))
		log.Debug().Uint32("seqno", master.SeqNo).Dur("took", time.Since(tm)).Msg("shards fetched")

		e := errgroup.Group{}

		for _, shard := range newShards {
			e.Go(func() error {
				log.Debug().Uint64("shard", uint64(shard.Shard)).Int32("wc", shard.Workchain).Msg("scanning shard")

				var block *tlbutils.Block

				err := v.retrier.Wrap(ctx, "fetch block", func() error {
					ctxBlock, err := v.api.Client().StickyContextNextNode(ctx)
					if err != nil {
						log.Debug().Err(err).Uint32("master", master.SeqNo).Int64("shard", shard.Shard).Msg("pick next node")
						return errors.Wrap(err, "pick next node")
					}

					waitCtx, cancel := context.WithTimeout(ctxBlock, defaultWaitNodeTimeout)
					defer cancel()

					block, err = v.api.WaitForBlock(master.SeqNo).GetBlockData(waitCtx, shard)
					if err != nil {
						log.Debug().Err(err).Uint32("master", master.SeqNo).Int64("shard", shard.Shard).Msg("get block")
						return errors.Wrap(err, "get block")
					}
					return nil
				})

				if err != nil || block == nil {
					return errors.Wrap(err, "fetch block")
				}

				err = func() error {
					shr := block.Extra.ShardAccountBlocks.BeginParse()
					shardAccBlocks, err := shr.LoadDict(256)
					if err != nil {
						return errors.Wrap(err, "load shard account blocks")
					}

					var wg sync.WaitGroup

					sab, err := shardAccBlocks.LoadAll()
					if err != nil {
						return errors.Wrap(err, "load all shard account blocks")
					}

					for _, kv := range sab {
						slc := kv.Value.MustToCell().BeginParse()
						if err = tlbutils.LoadFromCell(&tlbutils.CurrencyCollection{}, slc); err != nil {
							return errors.Wrap(err, "load aug currency collection of account block")
						}

						var ab tlbutils.AccountBlock
						if err = tlbutils.LoadFromCell(&ab, slc); err != nil {
							return errors.Wrap(err, "load account block")
						}

						allTx, err := ab.Transactions.LoadAll()
						if err != nil {
							return errors.Wrap(err, "load all transactions")
						}

						atomic.AddUint64(&transactionsNum, uint64(len(allTx)))

						for _, txKV := range allTx {
							slcTx := txKV.Value.MustToCell().BeginParse()
							if err = tlbutils.LoadFromCell(&tlbutils.CurrencyCollection{}, slcTx); err != nil {
								return errors.Wrap(err, "load aug currency collection of transaction")
							}

							var tx tlbutils.Transaction
							if err = tlbutils.LoadFromCell(&tx, slcTx.MustLoadRef()); err != nil {
								return errors.Wrap(err, "load transaction")
							}

							wg.Add(1)
							v.taskPool <- accFetchTask{
								master:   master,
								shard:    shard,
								tx:       &tx,
								addr:     addressutils.NewAddress(0, byte(shard.Workchain), ab.Addr),
								callback: wg.Done,
							}
							break
						}
					}

					log.Debug().
						Uint32("seqno", shard.SeqNo).Uint64("shard", uint64(shard.Shard)).Int32("wc", shard.Workchain).
						Int("affected_accounts", len(sab)).Uint64("transactions", transactionsNum).
						Msg("scanning transactions")

					wg.Wait()
					return nil
				}()

				if err != nil {
					log.Error().
						Uint32("seqno", shard.SeqNo).Uint64("shard", uint64(shard.Shard)).Int32("wc", shard.Workchain).
						Msg("failed to parse block, skipping. Fix issue and rescan later")
				}
				return nil
			})
		}

		if err := e.Wait(); err != nil {
			log.Error().Err(err).Msg("scan shard")
		}
		return transactionsNum, shardBlocksNum
	}
}
