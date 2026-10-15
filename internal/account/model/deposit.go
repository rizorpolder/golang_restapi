package model

type BalanceOperationType string

const (
	BalanceOperationDeposit BalanceOperationType = "DEPOSIT"
	BalanceOperationCredit  BalanceOperationType = "CREDIT"
)

type GetBalanceResponse struct {
	UserID  uint64
	Balance float32
}

type UpdateBalanceRequest struct {
	UserID  uint64
	Amount  float32
	Balance float32
	Type    BalanceOperationType
}
type UpdateBalanceResponse struct {
	UserID        uint64
	OldBalance    float32
	NewBalance    float32
	Amount        float32
	OperationType BalanceOperationType
}
