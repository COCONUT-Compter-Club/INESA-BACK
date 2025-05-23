package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type transactionServiceImpl struct {
	TransactionRepo repository.TransaksiRepo
	DB              *sql.DB
}

func NewTransactionService(transactionRepo repository.TransaksiRepo, db *sql.DB) TransactionService {
	return &transactionServiceImpl{
		TransactionRepo: transactionRepo,
		DB:              db,
	}
}

// GetAllTransaction implements TransactionService.
func (t *transactionServiceImpl) GetAllTransaction(ctx context.Context) ([]dto.TransactionResponse, error) {
	tx, err := t.DB.Begin()
	if err != nil {
		return []dto.TransactionResponse{}, fmt.Errorf("gagal memulai transaksi")
	}
	defer util.CommitOrRollBack(tx)

	transaction, err := t.TransactionRepo.GetAllTransaction(ctx, tx)
	if err != nil {
		return []dto.TransactionResponse{}, fmt.Errorf("gagal mendapatkan semua transaksi")
	}

	return util.ConvertTransactionToListResponseDTO(transaction), nil
}

// GetLastTransaction implements TransactionService.
func (t *transactionServiceImpl) GetLastTransaction(ctx context.Context) ([]dto.TransactionResponse, error) {
	tx, err := t.DB.Begin()
	if err != nil {
		panic(err)
	}
	defer util.CommitOrRollBack(tx)

	transaction, err := t.TransactionRepo.GetLastTransaction(ctx, tx)
	if err != nil {
		return []dto.TransactionResponse{}, fmt.Errorf("gagal mendapatkan transaksi terakhir")
	}

	return util.ConvertTransactionToListResponseDTO(transaction), nil
}
