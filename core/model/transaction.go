package model

import (
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/valyala/fastjson"
)

type Transaction struct {
	// Transaction identifiers
	AccountAddr string // Transaction identifier (AccountAddr or LT)
	LT          int64  // Logical time
	PrevTxHash  string // Previous transaction hash
	PrevTxLT    int64  // Previous transaction logical time

	// Address information
	Sender         string // Source address
	Receiver       string // Destination address
	SenderIsOurs   bool   // Flag indicating if sender is our address
	ReceiverIsOurs bool   // Flag indicating if receiver is our address

	// Financial information
	Amount    float64 // Transaction amount in TON
	TotalFees float64 // Total fees in TON
	ExitCode  int     // Compute phase exit code
	Success   bool    // Transaction execution success flag

	// Message information
	MessageType string // Type of message (INTERNAL, EXTERNAL_IN, etc.)
	Bounce      bool   // Bounce flag
	Bounced     bool   // Bounced flag
	Body        string // Transaction body (payload)

	// State information
	BlockID       string    // Block ID containing this transaction
	CreatedAt     time.Time // Transaction creation timestamp
	AccountStatus string    // Account status after transaction (ACTIVE, FROZEN, etc.)

	// Extra info
	ComputeGasUsed int    // Gas used for computation
	Description    string // Human-readable transaction description (optional)

}

func UnmarshalTransaction(data []byte) (*Transaction, error) {
	var p fastjson.Parser

	v, err := p.Parse(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "parse json")
	}

	tx := Transaction{}

	// Transaction identifiers
	if accountAddr := v.GetStringBytes("AccountAddr"); accountAddr != nil {
		tx.AccountAddr = string(accountAddr)
	}

	if lt := v.GetInt64("LT"); lt != 0 {
		tx.LT = lt
	}

	if prevTxHash := v.GetStringBytes("PrevTxHash"); prevTxHash != nil {
		tx.PrevTxHash = string(prevTxHash)
	}

	if prevLT := v.GetInt64("PrevTxLT"); prevLT != 0 {
		tx.PrevTxLT = prevLT
	}

	// IO information - handle messages
	if io := v.Get("IO"); io != nil {
		// Handle incoming message
		if inMsg := io.Get("In", "Msg"); inMsg != nil {
			// Address information
			if src := inMsg.Get("SrcAddr"); src != nil && src.Type() == fastjson.TypeString {
				tx.Sender = string(src.GetStringBytes())
			}

			if dst := inMsg.Get("DstAddr"); dst != nil && dst.Type() == fastjson.TypeString {
				tx.Receiver = string(dst.GetStringBytes())
			}

			// Message information
			if msgType := inMsg.Get("MsgType"); msgType != nil && msgType.Type() == fastjson.TypeString {
				tx.MessageType = string(msgType.GetStringBytes())
			}

			if bounce := inMsg.Get("Bounce"); bounce != nil {
				tx.Bounce = bounce.GetBool()
			}

			if bounced := inMsg.Get("Bounced"); bounced != nil {
				tx.Bounced = bounced.GetBool()
			}

			// Amount
			if amount := inMsg.Get("Amount"); amount != nil && amount.Type() == fastjson.TypeString {
				amountStr := string(amount.GetStringBytes())
				var amountInt float64
				if _, err := fmt.Sscanf(amountStr, "%f", &amountInt); err == nil {
					tx.Amount = amountInt / 1_000_000_000 // Convert from nano TON to TON
				}
			}

			// Body
			if body := inMsg.Get("Body"); body != nil && body.Type() == fastjson.TypeString {
				tx.Body = string(body.GetStringBytes())
			}
		}
	}

	// Time information
	if now := v.GetInt64("Now"); now != 0 {
		tx.CreatedAt = time.Unix(now, 0)
	}

	// State information
	if origStatus := v.GetStringBytes("OrigStatus"); origStatus != nil {
		// This is the starting status, not the ending one
	}

	if endStatus := v.GetStringBytes("EndStatus"); endStatus != nil {
		tx.AccountStatus = string(endStatus)
	}

	// Fee information
	if totalFees := v.Get("TotalFees", "Coins"); totalFees != nil && totalFees.Type() == fastjson.TypeString {
		feesStr := string(totalFees.GetStringBytes())
		var fees float64
		if _, err := fmt.Sscanf(feesStr, "%f", &fees); err == nil {
			tx.TotalFees = fees / 1_000_000_000 // Convert from nano TON to TON
		}
	}

	// Description contains compute phase and other details
	if desc := v.Get("Description"); desc != nil {
		// Extract success status
		if computePhase := desc.Get("ComputePhase", "Phase"); computePhase != nil {
			if success := computePhase.Get("Success"); success != nil {
				tx.Success = success.GetBool()
			}

			// Extract exit code and gas used
			if details := computePhase.Get("Details"); details != nil {
				tx.ExitCode = details.GetInt("ExitCode")
				tx.ComputeGasUsed = details.GetInt("GasUsed")
			}
		}
	}

	// Generate a description if success is true
	if tx.Success {
		tx.Description = fmt.Sprintf("Successfully transferred %.9f TON", tx.Amount)
	} else {
		tx.Description = fmt.Sprintf("Failed transaction with exit code %d", tx.ExitCode)
	}
	return &tx, nil
}
