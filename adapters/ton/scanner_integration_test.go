package ton_test

import (
	"context"
	"testing"
	"time"

	"github.com/kriuchkov/tonbeacon/adapters/ton"
	"github.com/stretchr/testify/suite"
	liteclientutils "github.com/xssnick/tonutils-go/liteclient"
	tonutils "github.com/xssnick/tonutils-go/ton"
)

type ScannerTestSuite struct {
	suite.Suite
	scanner *ton.Scanner
}

func (suite *ScannerTestSuite) SetupTest() {
	ctx := context.Background()

	client := liteclientutils.NewConnectionPool()

	err := client.AddConnectionsFromConfigUrl(ctx, testConfigUrl)
	suite.Require().NoError(err)

	liteClient := tonutils.NewAPIClient(client, tonutils.ProofCheckPolicyFast).WithRetry()
	suite.scanner = ton.NewScanner(liteClient, &ton.OptionsScanner{
		NumWorkers: 40,
	})
}

func (suite *ScannerTestSuite) TestRun() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	// Create a channel to receive scanner events
	resultsCh := make(chan interface{}, 100) // Adjust type based on your Scanner implementation

	// Start scanner with the channel
	suite.scanner.RunAsync(ctx, resultsCh)

	// Process results from the channel
	for {
		select {
		case result := <-resultsCh:
			suite.T().Logf("Received: %+v", result)
			// Add assertions if needed

		case <-ctx.Done():
			suite.T().Log("Scanner test completed")
			return
		}
	}
}

func TestScannerTestSuite(t *testing.T) {
	suite.Run(t, new(ScannerTestSuite))
}
