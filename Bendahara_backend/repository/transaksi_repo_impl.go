package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type transactionRepoImpl struct {
}

func NewTransactionRepo() TransaksiRepo {
	return &transactionRepoImpl{}
}

// GetAllTransaction implements TransaksiRepo.
func (t *transactionRepoImpl) GetAllTransaction(ctx context.Context, tx *sql.Tx) ([]model.Transaction, error) {
	query := "SELECT id_transaksi, tanggal, keterangan, jenis_transaksi, nominal FROM history_transaksi ORDER BY tanggal DESC"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil transaksi: %v", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		transaction := model.Transaction{}
		var tanggalStr string // Simpan tanggal sebagai string terlebih dahulu

		// Scan ke variabel sementara (tanggalStr)
		err := rows.Scan(&transaction.Id, &tanggalStr, &transaction.Keterangan, &transaction.JenisTransaksi, &transaction.Nominal)
		if err != nil {
			return nil, fmt.Errorf("gagal memindai transaksi: %v", err)
		}

		// Parsing string ke time.Time
		transaction.Tanggal, err = time.Parse("2006-01-02 15:04:05", tanggalStr) // Sesuaikan format dengan data di database
		if err != nil {
			return nil, fmt.Errorf("gagal menguraikan tanggal: %v", err)
		}

		transactions = append(transactions, transaction)
	}

	// Periksa error setelah iterasi
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("kesalahan setelah melakukan iterasi baris: %v", err)
	}

	return transactions, nil
}

// GetLastTransaction implements TransaksiRepo.
func (t *transactionRepoImpl) GetLastTransaction(ctx context.Context, tx *sql.Tx) ([]model.Transaction, error) {
	query := "SELECT id_transaksi, tanggal, keterangan, jenis_transaksi, nominal FROM history_transaksi ORDER BY tanggal DESC LIMIT 5"
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil transaksi terakhir: %v", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		transaction := model.Transaction{}
		var tanggalStr string // Simpan tanggal sebagai string terlebih dahulu

		// Scan ke variabel sementara (tanggalStr)
		err := rows.Scan(&transaction.Id, &tanggalStr, &transaction.Keterangan, &transaction.JenisTransaksi, &transaction.Nominal)
		if err != nil {
			return nil, fmt.Errorf("gagal memindai transaksi: %v", err)
		}

		// Parsing string ke time.Time
		transaction.Tanggal, err = time.Parse("2006-01-02 15:04:05", tanggalStr) // Sesuaikan format dengan data di database
		if err != nil {
			return nil, fmt.Errorf("gagal mengurai tanggal: %v", err)
		}

		transactions = append(transactions, transaction)
	}

	// Periksa error setelah iterasi
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("kesalahan setelah melakukan iterasi baris: %v", err)
	}

	return transactions, nil
}
