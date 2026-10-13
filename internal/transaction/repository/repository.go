package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"golang_restapi/internal/transaction/model"
	"golang_restapi/internal/transaction/repository/mapper"
	repmodel "golang_restapi/internal/transaction/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewRepository(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) GetTransactions(ctx context.Context, params model.GetTransactionParams) ([]model.Transaction, error) {
	var transactions []repmodel.Transaction
	query := r.db.WithContext(ctx).Model(&repmodel.Transaction{})

	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if params.Type != nil {
		query = query.Where("type = ?", *params.Type)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.DateFrom != nil {
		query = query.Where("created_at >= ?", *params.DateFrom)
	}
	if params.DateTo != nil {
		query = query.Where("created_at <= ?", *params.DateTo)
	}

	res := query.Offset(params.Offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&transactions)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to fetch transactions")
		return nil, fmt.Errorf("failed to get transactions: %w", res.Error)
	}
	return mapper.RepoTransactionsToTransactions(transactions), nil
}

func (r *Repository) GetTransactionDetails(ctx context.Context, transactionId uint64) (model.TransactionDetails, error) {
	var transaction repmodel.Transaction
	res := r.db.WithContext(ctx).
		Model(&repmodel.Transaction{}).
		Where("id = ?", transactionId).
		First(&transaction)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.TransactionDetails{}, fmt.Errorf("transaction %d not found", transactionId)
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to fetch transaction details")
		return model.TransactionDetails{}, fmt.Errorf("failed to fetch transaction details: %w", res.Error)
	}

	var entries []repmodel.TransactionEntry
	res = r.db.WithContext(ctx).Model(&repmodel.TransactionEntry{}).
		Where("transaction_id =?", transactionId).
		Order("created_at ASC").
		Find(&entries)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to fetch transaction details")
		return model.TransactionDetails{}, fmt.Errorf("failed to fetch transaction details: %w", res.Error)
	}
	return model.TransactionDetails{
		Transaction: mapper.RepoTransactionToTransaction(transaction),
		Entries:     mapper.RepoTransactionEntriesToTransactionEntries(entries),
	}, nil
}

func (r *Repository) Withdraw(ctx context.Context, params model.WithdrawParams) (model.TransactionDetails, error) {
	var result model.TransactionDetails
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transaction := model.Transaction{
			UserID: params.AccountID,
			Amount: int64(params.Amount * 100),
			Status: model.TransactionStatusPending,
			Type:   model.TransactionTypeWithdraw,
		}

		transactionRepo := mapper.TransactionToRepoTransaction(transaction)
		res := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&transaction)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed to create withdraw transaction")
			return fmt.Errorf("failed to create withdraw transaction: %w", res.Error)
		}

		entry := model.TransactionEntry{
			TransactionID: transactionRepo.ID,
			AccountID:     params.AccountID,
			Direction:     model.TransactionEntryDirectionCredit,
			Amount:        params.Amount,
		}

		entryRepo := mapper.TransactionEntryToRepoTransactionEntry(entry)
		res = tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&entry)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed to create withdraw transaction entry")
			return fmt.Errorf("failed to create withdraw transaction entry: %w", res.Error)
		}

		updateData := mapper.UpdateTransactionToRepoTransaction(model.UpdateTransaction{Status: model.TransactionStatusCompleted})
		res = tx.Where("id =?", transactionRepo.ID).Updates(&updateData)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed to update withdraw transaction entry")
			return fmt.Errorf("failed to update withdraw transaction entry: %w", res.Error)
		}

		transactionRepo.Status = string(model.TransactionStatusCompleted)
		result = model.TransactionDetails{
			Transaction: mapper.RepoTransactionToTransaction(transactionRepo),
			Entries: []model.TransactionEntry{
				mapper.RepoTransactionEntryToTransactionEntry(entryRepo),
			},
		}
		return nil
	})
	if err != nil {
		return model.TransactionDetails{}, fmt.Errorf("failed to create withdraw transaction: %w", err)
	}

	return result, nil
}

func (r *Repository) Deposit(ctx context.Context, params model.DepositParams) (model.TransactionDetails, error) {
	var result model.TransactionDetails
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transaction := model.Transaction{
			UserID: params.UserID,
			Amount: int64(params.Amount * 100),
			Status: model.TransactionStatusPending,
			Type:   model.TransactionTypeDeposit,
		}

		transactionRepo := mapper.TransactionToRepoTransaction(transaction)
		res := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&transactionRepo)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed create transaction deposit")
			return fmt.Errorf("failed to create transaction: %w", res.Error)
		}

		entry := model.TransactionEntry{
			TransactionID: transaction.ID,
			AccountID:     params.UserID,
			Direction:     model.TransactionEntryDirectionDebit,
			Amount:        params.Amount,
		}

		entryRepo := mapper.TransactionEntryToRepoTransactionEntry(entry)
		res = tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&entryRepo)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed create transaction entry in deposit")
			return fmt.Errorf("failed to create transaction entry: %w", res.Error)
		}

		updateData := mapper.UpdateTransactionToRepoTransaction(model.UpdateTransaction{
			Status: model.TransactionStatusCompleted,
		})

		res = tx.Where("id = ?", transaction.ID).Updates(&updateData)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed update transaction entry in deposit")
			return fmt.Errorf("failed to update transaction entry: %w", res.Error)
		}

		transactionRepo.Status = string(model.TransactionStatusCompleted)
		result = model.TransactionDetails{
			Transaction: mapper.RepoTransactionToTransaction(transactionRepo),
			Entries:     []model.TransactionEntry{mapper.RepoTransactionEntryToTransactionEntry(entryRepo)},
		}
		return nil
	})
	if err != nil {
		return model.TransactionDetails{}, err
	}
	return result, nil
}
func (r *Repository) Transfer(ctx context.Context, params model.TransferParams) (model.TransactionDetails, error) {
	var result model.TransactionDetails
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transaction := model.Transaction{
			UserID: params.UserID,
			Amount: int64(params.Amount * 100),
			Status: model.TransactionStatusPending,
			Type:   model.TransactionTypeTransfer,
		}
		transactionRepo := mapper.TransactionToRepoTransaction(transaction)
		res := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&transactionRepo)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed create transaction transfer")
			return fmt.Errorf("failed to create transaction: %w", res.Error)
		}

		creditEntry := model.TransactionEntry{
			TransactionID: transactionRepo.ID,
			AccountID:     params.UserID,
			Direction:     model.TransactionEntryDirectionCredit,
			Amount:        params.Amount,
		}
		creditEntryRepo := mapper.TransactionEntryToRepoTransactionEntry(creditEntry)
		res = tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&creditEntryRepo)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("ailed to create credit transaction entry")
			return fmt.Errorf("failed to create credit transaction entry: %w", res.Error)
		}

		debitEntry := model.TransactionEntry{
			TransactionID: transactionRepo.ID,
			AccountID:     params.Recipient,
			Direction:     model.TransactionEntryDirectionDebit,
			Amount:        params.Amount,
		}

		debitEntryRepo := mapper.TransactionEntryToRepoTransactionEntry(debitEntry)
		res = tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&debitEntryRepo)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("ailed to create debit transaction entry")
			return fmt.Errorf("failed to create debit transaction entry: %w", res.Error)
		}

		updateData := mapper.UpdateTransactionToRepoTransaction(model.UpdateTransaction{
			Status: model.TransactionStatusCompleted,
		})
		res = tx.Where("id = ?", transactionRepo.ID).Updates(&updateData)
		if res.Error != nil {
			r.logger.Err(res.Error).Msg("failed update transaction entry in transfer")
			return fmt.Errorf("failed to update transaction entry: %w", res.Error)
		}
		transactionRepo.Status = string(model.TransactionStatusCompleted)
		result = model.TransactionDetails{
			Transaction: mapper.RepoTransactionToTransaction(transactionRepo),
			Entries: []model.TransactionEntry{
				mapper.RepoTransactionEntryToTransactionEntry(creditEntryRepo),
				mapper.RepoTransactionEntryToTransactionEntry(debitEntryRepo),
			},
		}
		return nil
	})
	if err != nil {
		return model.TransactionDetails{}, err
	}
	return result, nil
}
