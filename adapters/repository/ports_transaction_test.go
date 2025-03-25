package repository

import (
	"context"
	"time"

	"github.com/kriuchkov/tonbeacon/core/model"
)

func (suite *RepositoryTestSuite) TestInsertTransaction() {
	testTransaction := &model.Transaction{
		AccountAddr:    "EQBAGi6wUF6SvjQtWyP5OvniQfYSI3Q-eTvM4BLiJztcRahv",
		LT:             1674235553000,
		PrevTxHash:     "97b7bf0154d3b1a3ce9ac692944e53f518f819b7e09ba567a61ad0bd1724fc30",
		PrevTxLT:       1674235552981,
		Sender:         "EQDrLq-X6jKZNHAScgghh0h1iijd-NQO0lSPEoZ4exDtOFLt",
		Receiver:       "EQBAGi6wUF6SvjQtWyP5OvniQfYSI3Q-eTvM4BLiJztcRahv",
		SenderIsOurs:   false,
		ReceiverIsOurs: true,
		Amount:         2.75,
		TotalFees:      0.01,
		ExitCode:       0,
		Success:        true,
		MessageType:    "INTERNAL",
		Bounce:         true,
		Bounced:        false,
		Body:           "te6cckEBAQEADgAAGKzAUQkAAAAAAAAAAADA84rE",
		BlockID:        "(-1,8000001,12345)",
		CreatedAt:      time.Now(),
		AccountStatus:  "ACTIVE",
		ComputeGasUsed: 10123,
		Description:    "Incoming TON transfer",
	}

	ctx := context.Background()
	inserted, err := suite.adapter.InsertTransaction(ctx, testTransaction)

	suite.NoError(err)
	suite.NotNil(inserted)

	suite.Equal(testTransaction.AccountAddr, inserted.AccountAddr)
	suite.Equal(testTransaction.LT, inserted.LT)
	suite.Equal(testTransaction.PrevTxHash, inserted.PrevTxHash)
	suite.Equal(testTransaction.PrevTxLT, inserted.PrevTxLT)
	suite.Equal(testTransaction.Sender, inserted.Sender)
	suite.Equal(testTransaction.Receiver, inserted.Receiver)
	suite.Equal(testTransaction.Amount, inserted.Amount)
	suite.Equal(testTransaction.TotalFees, inserted.TotalFees)
}
