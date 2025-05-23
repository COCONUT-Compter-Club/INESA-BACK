package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/repository"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type laporanKeuanganServiceImpl struct {
	LaporanRepo repository.LaporanKeuanganRepo
	DB          *sql.DB
}

func NewLaporanKeuanganService(laporanRepo repository.LaporanKeuanganRepo, db *sql.DB) LaporanKeuanganService {
	return &laporanKeuanganServiceImpl{
		LaporanRepo: laporanRepo,
		DB:          db,
	}
}

// GetAllLaporan implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetAllLaporan(ctx context.Context) ([]dto.LaporanKeuanganResponse, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, fmt.Errorf("gagal memulai transaksi")
	}
	defer util.CommitOrRollBack(tx)

	laporan, err := l.LaporanRepo.GetAllLaporan(ctx, tx)
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, fmt.Errorf("gagal mendapatkan laporan keuangan")
	}

	return util.ConvertLaporanToListResponseDTO(laporan), nil
}

// GetLastBalance implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetLastBalance(ctx context.Context) (int64, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("gagal memulai transaksi")
	}
	defer util.CommitOrRollBack(tx)

	saldo, err := l.LaporanRepo.GetLastBalance(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("gagal mendapatkan saldo terakhir: %v", err)
	}

	return saldo, nil
}

// GetTotalExpenditure implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetTotalExpenditure(ctx context.Context) (uint64, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("gagal memulai transaksi")
	}
	defer util.CommitOrRollBack(tx)

	laporan, err := l.LaporanRepo.GetAllLaporan(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("gagal mendapatkan laporan keuangan")
	}

	var totalPengeluaran uint64
	for _, data := range laporan {
		totalPengeluaran += data.Pengeluaran
	}

	return totalPengeluaran, nil
}

// GetTotalIncome implements LaporanKeuanganService.
func (l *laporanKeuanganServiceImpl) GetTotalIncome(ctx context.Context) (uint64, error) {
	tx, err := l.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("gagal memulai transaksi")
	}
	defer util.CommitOrRollBack(tx)

	laporan, err := l.LaporanRepo.GetAllLaporan(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("gagal mendapatkan laporan keuangan")
	}

	var totalPemasukan uint64
	for _, data := range laporan {
		totalPemasukan += data.Pemasukan
	}

	return totalPemasukan, nil
}

func (l *laporanKeuanganServiceImpl) GetLaporanByDateRange(ctx context.Context, startDate string, endDate string) ([]dto.LaporanKeuanganResponse, error) {
	// Mulai transaksi
	tx, err := l.DB.Begin()
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, fmt.Errorf("gagal memulai transaksi")
	}
	defer util.CommitOrRollBack(tx)

	// Panggil repository untuk mendapatkan data laporan berdasarkan rentang tanggal
	laporans, err := l.LaporanRepo.GetLaporanByDateRange(ctx, tx, startDate, endDate)
	if err != nil {
		return []dto.LaporanKeuanganResponse{}, fmt.Errorf("gagal mendapatkan laporan berdasarkan rentang tanggal: %v", err)
	}

	// Konversi data laporan ke DTO
	return util.ConvertLaporanToListResponseDTO(laporans), nil
}