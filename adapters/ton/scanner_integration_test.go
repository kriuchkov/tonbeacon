package ton_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	liteclientutils "github.com/xssnick/tonutils-go/liteclient"
	tonutils "github.com/xssnick/tonutils-go/ton"

	"github.com/kriuchkov/tonbeacon/adapters/ton"
)

type ScannerTestSuite struct {
	suite.Suite
	scanner *ton.Scanner
}

func (suite *ScannerTestSuite) SetupTest() {
	ctx := context.Background()

	client := liteclientutils.NewConnectionPool()

	err := client.AddConnectionsFromConfigUrl(ctx, testConfigURL)
	suite.Require().NoError(err)

	liteClient := tonutils.NewAPIClient(client, tonutils.ProofCheckPolicyFast).WithRetry()
	suite.scanner = ton.NewScanner(liteClient, &ton.OptionsScanner{
		NumWorkers: 40,
	})
}

func (suite *ScannerTestSuite) TestRun() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	resultsCh := make(chan any, 100)
	suite.scanner.RunAsync(ctx, resultsCh)

	for {
		select {
		case result := <-resultsCh:
			suite.T().Logf("Received: %+v", result)

		case <-ctx.Done():
			suite.T().Log("Scanner test completed")
			return
		}
	}
}

func TestScannerTestSuite(t *testing.T) {
	suite.Run(t, new(ScannerTestSuite))
}
