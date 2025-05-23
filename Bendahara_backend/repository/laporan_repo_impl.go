package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

type laporanKeuanganRepoImpl struct {}

// NewLaporanKeuanganRepo creates a new instance of the laporan keuangan repository
func NewLaporanKeuanganRepo() LaporanKeuanganRepo {
	return &laporanKeuanganRepoImpl{}
}

// GetAllLaporan retrieves all laporan keuangan, sorted by tanggal descending
func (l *laporanKeuanganRepoImpl) GetAllLaporan(ctx context.Context, tx *sql.Tx) ([]model.LaporanKeuangan, error) {
	var laporans []model.LaporanKeuangan
	// Include column `nota` in the SELECT
	query := `SELECT id_laporan, tanggal, keterangan, nota, pemasukan, pengeluaran, saldo
		FROM laporan_keuangan
		ORDER BY tanggal DESC`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return laporans, err
	}
	defer rows.Close()

	for rows.Next() {
		var laporan model.LaporanKeuangan
		var tanggalRaw interface{}

		// Scan including the `nota` field
		err := rows.Scan(
			&laporan.Id,
			&tanggalRaw,
			&laporan.Keterangan,
			&laporan.Nota,
			&laporan.Pemasukan,
			&laporan.Pengeluaran,
			&laporan.Saldo,
		)
		if err != nil {
			return laporans, err
		}

		// Convert tanggal to time.Time
		switch v := tanggalRaw.(type) {
		case time.Time:
			laporan.Tanggal = v
		case []byte:
			parsed, perr := time.Parse("2006-01-02 15:04:05", string(v))
			if perr != nil {
				return laporans, fmt.Errorf("failed to parse tanggal: %v", perr)
			}
			laporan.Tanggal = parsed
		default:
			return laporans, fmt.Errorf("unsupported type for tanggal: %T", v)
		}

		laporans = append(laporans, laporan)
	}

	if err = rows.Err(); err != nil {
		return laporans, err
	}

	return laporans, nil
}

// GetLastBalance retrieves the most recent saldo
func (l *laporanKeuanganRepoImpl) GetLastBalance(ctx context.Context, tx *sql.Tx) (int64, error) {
	var saldo int64
	query := `SELECT saldo FROM laporan_keuangan ORDER BY tanggal DESC LIMIT 1`
	err := tx.QueryRowContext(ctx, query).Scan(&saldo)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch previous saldo: %v", err)
	}
	return saldo, nil
}

// GetLaporanByDateRange retrieves laporan keuangan within a given date range
func (l *laporanKeuanganRepoImpl) GetLaporanByDateRange(ctx context.Context, tx *sql.Tx, startDate, endDate string) ([]model.LaporanKeuangan, error) {
	var laporans []model.LaporanKeuangan
	// Include `nota` in SELECT
	query := `SELECT id_laporan, tanggal, keterangan, nota, pemasukan, pengeluaran, saldo
		FROM laporan_keuangan
		WHERE tanggal BETWEEN ? AND ?
		ORDER BY tanggal ASC`

	rows, err := tx.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return laporans, fmt.Errorf("failed to fetch laporan by date range: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var laporan model.LaporanKeuangan
		var tanggalStr string

		err := rows.Scan(
			&laporan.Id,
			&tanggalStr,
			&laporan.Keterangan,
			&laporan.Nota,
			&laporan.Pemasukan,
			&laporan.Pengeluaran,
			&laporan.Saldo,
		)
		if err != nil {
			return laporans, fmt.Errorf("failed to scan laporan: %v", err)
		}

		// Parse tanggal string
		parsed, perr := time.Parse("2006-01-02 15:04:05", tanggalStr)
		if perr != nil {
			return laporans, fmt.Errorf("failed to parse tanggal: %v", perr)
		}
		laporan.Tanggal = parsed

		laporans = append(laporans, laporan)
	}

	if err = rows.Err(); err != nil {
		return laporans, fmt.Errorf("error after iterating rows: %v", err)
	}

	return laporans, nil
}
