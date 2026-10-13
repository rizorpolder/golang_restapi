package mapper

import (
	"golang_restapi/internal/transaction/model"
	repmodel "golang_restapi/internal/transaction/repository/model"
)

func RepoTransactionsToTransactions(transaction []repmodel.Transaction) []model.Transaction {
	result := make([]model.Transaction, len(transaction))

	for i, transaction := range transaction {
		result[i] = RepoTransactionToTransaction(transaction)
	}
	return result
}

func TransactionsToRepoTransactions(transaction []model.Transaction) []repmodel.Transaction {
	result := make([]repmodel.Transaction, len(transaction))

	for i, transaction := range transaction {
		result[i] = TransactionToRepoTransaction(transaction)
	}
	return result
}

func RepoTransactionToTransaction(transaction repmodel.Transaction) model.Transaction {
	return model.Transaction{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Amount:    transaction.Amount,
		Status:    model.TransactionStatus(transaction.Status),
		Type:      model.TransactionType(transaction.Type),
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}
}

func TransactionToRepoTransaction(transaction model.Transaction) repmodel.Transaction {
	return repmodel.Transaction{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Amount:    transaction.Amount,
		Status:    string(transaction.Status),
		Type:      string(transaction.Type),
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}
}

func RepoTransactionEntriesToTransactionEntries(entries []repmodel.TransactionEntry) []model.TransactionEntry {
	result := make([]model.TransactionEntry, len(entries))
	for i, entry := range entries {
		result[i] = RepoTransactionEntryToTransactionEntry(entry)
	}
	return result
}

func TransactionEntriesToRepoTransactionEntries(entries []model.TransactionEntry) []repmodel.TransactionEntry {
	result := make([]repmodel.TransactionEntry, len(entries))
	for i, entry := range entries {
		result[i] = TransactionEntryToRepoTransactionEntry(entry)
	}
	return result
}

func RepoTransactionEntryToTransactionEntry(entry repmodel.TransactionEntry) model.TransactionEntry {
	return model.TransactionEntry{
		ID:            entry.ID,
		TransactionID: entry.TransactionID,
		AccountID:     entry.AccountID,
		Amount:        entry.Amount,
		Direction:     model.TransactionEntryDirection(entry.Direction),
		CreatedAt:     entry.CreatedAt,
		UpdatedAt:     entry.UpdatedAt,
	}
}

func TransactionEntryToRepoTransactionEntry(entry model.TransactionEntry) repmodel.TransactionEntry {
	return repmodel.TransactionEntry{
		ID:            entry.ID,
		TransactionID: entry.TransactionID,
		AccountID:     entry.AccountID,
		Amount:        entry.Amount,
		Direction:     string(entry.Direction),
		CreatedAt:     entry.CreatedAt,
		UpdatedAt:     entry.UpdatedAt,
	}
}

func UpdateTransactionToRepoTransaction(transaction model.UpdateTransaction) repmodel.UpdateTransaction {
	return repmodel.UpdateTransaction{
		Status: string(transaction.Status),
	}
}
