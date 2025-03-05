package ton

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"slices"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	addressutils "github.com/xssnick/tonutils-go/address"
	tlbutils "github.com/xssnick/tonutils-go/tlb"
	tonutils "github.com/xssnick/tonutils-go/ton"
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
	for y := 0; y < threads; y++ {
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
) (ret []*tonutils.BlockIDExt, lastTime time.Time, err error) {
	if shard.Workchain != 0 {
		return nil, time.Time{}, nil
	}

	if slices.ContainsFunc(prevShards, shard.Equals) {
		return nil, time.Time{}, nil
	}

	b, err := api.GetBlockData(ctx, shard)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("get block data: %w", err)
	}

	parents, err := b.BlockInfo.GetParentBlocks()
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("get parent blocks (%d:%x:%d): %w", shard.Workchain, uint64(shard.Shard), shard.Shard, err)
	}

	genTime := time.Unix(int64(b.BlockInfo.GenUtime), 0)

	for _, parent := range parents {
		ext, _, err := v.getNotSeenShards(ctx, api, parent, prevShards)
		if err != nil {
			return nil, time.Time{}, err
		}
		ret = append(ret, ext...)
	}
	return append(ret, shard), genTime, nil
}

func (v *Scanner) fetchBlock(ctx context.Context, master *tonutils.BlockIDExt) (transactionsNum, shardBlocksNum uint64) {
	log.Debug().Uint32("seqno", master.SeqNo).Msg("scanning master")

	tm := time.Now()
	for {
		select {
		case <-ctx.Done():
			log.Warn().Uint32("master", master.SeqNo).Msg("ctx done")
			return
		default:
		}

		prevMaster, err := v.api.WaitForBlock(master.SeqNo-1).LookupBlock(ctx, master.Workchain, master.Shard, master.SeqNo-1)
		if err != nil {
			log.Debug().Err(err).Uint32("seqno", master.SeqNo-1).Msg("failed to get prev master block")
			time.Sleep(300 * time.Millisecond)
			continue
		}

		prevShards, err := v.api.GetBlockShardsInfo(ctx, prevMaster)
		if err != nil {
			log.Debug().Err(err).Uint32("master", master.SeqNo).Msg("failed to get shards on block")
			time.Sleep(300 * time.Millisecond)
			continue
		}

		// getting information about other work-chains and shards of master block
		currentShards, err := v.api.GetBlockShardsInfo(ctx, master)
		if err != nil {
			log.Debug().Err(err).Uint32("master", master.SeqNo).Msg("failed to get shards on block")
			time.Sleep(300 * time.Millisecond)
			continue
		}
		log.Debug().Uint32("seqno", master.SeqNo).Dur("took", time.Since(tm)).Msg("shards fetched")

		// shards in master block may have holes, e.g. shard seqno 2756461, then 2756463, and no 2756462 in master chain
		// thus we need to scan a bit back in case of discovering a hole, till last seen, to fill the misses.
		var newShards []*tonutils.BlockIDExt
		for _, shard := range currentShards {
			for {
				select {
				case <-ctx.Done():
					log.Warn().Uint32("master", master.SeqNo).Msg("ctx done")
					return
				default:
				}

				notSeen, _, err := v.getNotSeenShards(ctx, v.api, shard, prevShards)
				if err != nil {
					log.Debug().Err(err).Uint32("master", master.SeqNo).Msg("failed to get not seen shards on block")
					time.Sleep(300 * time.Millisecond)
					continue
				}

				newShards = append(newShards, notSeen...)
				break
			}
		}

		log.Debug().Uint32("seqno", master.SeqNo).Dur("took", time.Since(tm)).Msg("not seen shards fetched")

		var shardsWg sync.WaitGroup
		shardsWg.Add(len(newShards))
		shardBlocksNum = uint64(len(newShards))
		// for each shard block getting transactions
		for _, shard := range newShards {
			log.Debug().Uint32("seqno", shard.SeqNo).Uint64("shard", uint64(shard.Shard)).Int32("wc", shard.Workchain).Msg("scanning shard")

			go func(shard *tonutils.BlockIDExt) {
				defer shardsWg.Done()

				var block *tlbutils.Block
				{
					ctx := ctx
					for z := 0; z < 20; z++ { // TODO: retry without loosing
						ctx, err = v.api.Client().StickyContextNextNode(ctx)
						if err != nil {
							log.Debug().Err(err).Uint32("master", master.SeqNo).Int64("shard", shard.Shard).
								Uint32("shard_seqno", shard.SeqNo).Msg("failed to pick next node")
							break
						}

						qCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
						block, err = v.api.WaitForBlock(master.SeqNo).GetBlockData(qCtx, shard)
						cancel()
						if err != nil {
							log.Debug().Err(err).Uint32("master", master.SeqNo).Int64("shard", shard.Shard).
								Uint32("shard_seqno", shard.SeqNo).Msg("failed to get block")
							time.Sleep(200 * time.Millisecond)
							continue
						}
						break
					}
				}

				err = func() error {
					if block == nil {
						return fmt.Errorf("failed to fetch block")
					}

					shr := block.Extra.ShardAccountBlocks.BeginParse()
					shardAccBlocks, err := shr.LoadDict(256)
					if err != nil {
						return fmt.Errorf("faled to load shard account blocks dict: %w", err)
					}

					var wg sync.WaitGroup

					sab := shardAccBlocks.All()
					for _, kv := range sab {
						slc := kv.Value.BeginParse()
						if err = tlbutils.LoadFromCell(&tlbutils.CurrencyCollection{}, slc); err != nil {
							return fmt.Errorf("faled to load aug currency collection of account block dict: %w", err)
						}

						var ab tlbutils.AccountBlock
						if err = tlbutils.LoadFromCell(&ab, slc); err != nil {
							return fmt.Errorf("faled to parse account block: %w", err)
						}

						allTx := ab.Transactions.All()
						transactionsNum += uint64(len(allTx))
						for _, txKV := range allTx {
							slcTx := txKV.Value.BeginParse()
							if err = tlbutils.LoadFromCell(&tlbutils.CurrencyCollection{}, slcTx); err != nil {
								return fmt.Errorf("faled to load aug currency collection of transactions dict: %w", err)
							}

							var tx tlbutils.Transaction
							if err = tlbutils.LoadFromCell(&tx, slcTx.MustLoadRef()); err != nil {
								return fmt.Errorf("faled to parse transaction: %w", err)
							}

							wg.Add(1)
							v.taskPool <- accFetchTask{
								master:   master,
								shard:    shard,
								tx:       &tx,
								addr:     addressutils.NewAddress(0, byte(shard.Workchain), ab.Addr),
								callback: wg.Done,
							}
							// 1 tx for account is enough for us, as a reference
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
			}(shard)
		}

		shardsWg.Wait()
		return
	}
}
