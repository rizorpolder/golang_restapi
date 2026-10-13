package model

import "time"

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETE"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCanceled  TransactionStatus = "CANCELED"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

type Transaction struct {
	ID        uint64
	UserID    uint64
	Amount    int64 //0.00
	Status    TransactionStatus
	Type      TransactionType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateTransaction struct {
	Status TransactionStatus
}
type TransactionEntryDirection string

const (
	TransactionEntryDirectionDebit  TransactionEntryDirection = "DEBIT"
	TransactionEntryDirectionCredit TransactionEntryDirection = "CREDIT"
)

type TransactionEntry struct {
	ID            uint64
	TransactionID uint64
	AccountID     uint64
	Direction     TransactionEntryDirection
	Amount        float64 //0.00
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateTransactionEntry struct {
	TransactionID uint64
	AccountID     uint64
	Direction     TransactionEntryDirection
	Amount        float64 //0.00
}

type TransactionDetails struct {
	Transaction Transaction
	Entries     []TransactionEntry
}

type GetTransactionParams struct {
	UserID   *uint64
	Type     *string
	Status   *string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

type DepositParams struct {
	UserID uint64
	Amount float64 //0.00
}

type WithdrawParams struct {
	AccountID uint64
	Amount    float64 //0.00
}

type TransferParams struct {
	UserID    uint64
	Recipient uint64
	Amount    float64 //0.00
}
