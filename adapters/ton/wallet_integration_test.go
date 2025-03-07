package ton_test

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
	liteclientutils "github.com/xssnick/tonutils-go/liteclient"
	tonutils "github.com/xssnick/tonutils-go/ton"
	walletutils "github.com/xssnick/tonutils-go/ton/wallet"

	tonadapter "github.com/kriuchkov/tonbeacon/adapters/ton"
)

var testConfigURL = "https://tonutils.com/testnet-global.config.json"

type WalletManagerTestSuite struct {
	suite.Suite
	masterWallet  *walletutils.Wallet
	walletAdapter *tonadapter.WalletAdapter
}

func (suite *WalletManagerTestSuite) SetupTest() {
	ctx := context.Background()

	client := liteclientutils.NewConnectionPool()

	err := client.AddConnectionsFromConfigUrl(ctx, testConfigURL)
	suite.Require().NoError(err)

	seed := walletutils.NewSeed()
	suite.Require().NotEmpty(seed)
	suite.T().Logf("Seed phrase: %s\n", seed)

	liteClient := tonutils.NewAPIClient(client, tonutils.ProofCheckPolicyFast).WithRetry()

	masterWallet, err := walletutils.FromSeed(liteClient, seed, walletutils.ConfigV5R1Final{
		NetworkGlobalID: walletutils.TestnetGlobalID,
		Workchain:       0,
	})
	suite.Require().NoError(err)

	suite.T().Logf("Master wallet address: %s\n", masterWallet.WalletAddress().String())
	suite.masterWallet = masterWallet

	suite.walletAdapter = tonadapter.NewWalletAdapter(liteClient, masterWallet)
}

func (suite *WalletManagerTestSuite) TestWalletManager() {
	ctx := context.Background()

	subWallet, err := suite.walletAdapter.CreateWallet(ctx, 1)
	suite.Require().NoError(err)

	suite.T().Logf("Subwallet address 1: %s\n", subWallet.WalletAddress().String())

	subWallet2, err := suite.walletAdapter.CreateWallet(ctx, math.MaxInt32-1)
	suite.Require().NoError(err)

	suite.T().Logf("Subwallet address 2: %s\n", subWallet2.WalletAddress().String())

	subWallet3, err := suite.walletAdapter.CreateWallet(ctx, math.MaxInt32)
	suite.Require().NoError(err)

	suite.T().Logf("Subwallet address 3: %s\n", subWallet3.WalletAddress().String())
}

func (suite *WalletManagerTestSuite) TestGetBalance() {
	ctx := context.Background()

	subWallet, err := suite.walletAdapter.CreateWallet(ctx, 1)
	suite.Require().NoError(err)

	suite.T().Logf("Subwallet address: %s\n", subWallet.WalletAddress().String())

	balance, err := suite.walletAdapter.GetBalance(ctx, 1)
	suite.Require().NoError(err)

	suite.T().Logf("Balance: %v\n", balance)
}

func TestWalletManagerTestSuite(t *testing.T) {
	suite.Run(t, new(WalletManagerTestSuite))
}
