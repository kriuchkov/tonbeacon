package ton_test

import (
	"context"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	walletutils "github.com/xssnick/tonutils-go/ton/wallet"

	tonadapter "github.com/kriuchkov/tonbeacon/adapters/ton"
	"github.com/kriuchkov/tonbeacon/core/consts"
	"github.com/kriuchkov/tonbeacon/pkg/common"
)

var (
	testConfigURL = consts.TestNetConfigURL
	testSeed      = os.Getenv("TON_SEED")
	testVersion   = walletutils.V4R2
)

type WalletManagerTestSuite struct {
	suite.Suite
	masterWallet  *walletutils.Wallet
	walletAdapter *tonadapter.WalletAdapter
}

func (suite *WalletManagerTestSuite) SetupTest() {
	ctx := context.Background()

	liteClient, err := common.SetupLiteClient(ctx, false)
	suite.Require().NoError(err)

	if testSeed == "" {
		suite.T().Fatal("TON_SEED env variable is not set")
	}

	masterWallet, err := walletutils.FromSeed(liteClient, strings.Split(testSeed, " "), testVersion)
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

	balance, err := suite.walletAdapter.GetBalance(ctx, 0)
	suite.Require().NoError(err)

	suite.T().Logf("Balance: %v\n", balance)
}

func TestWalletManagerTestSuite(t *testing.T) {
	suite.Run(t, new(WalletManagerTestSuite))
}
