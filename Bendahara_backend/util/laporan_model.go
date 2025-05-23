package util

import (
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/model"
)

func ConvertLaporanToResponseDTO(laporan model.LaporanKeuangan) dto.LaporanKeuanganResponse {
	// Format tanggal sesuai dengan format yang diinginkan (misalnya: "02-01-2006 15:04")
	formattedTanggal := laporan.Tanggal.Format("02-01-2006 15:04")

	return dto.LaporanKeuanganResponse{
		Id:          laporan.Id,
		Tanggal:     formattedTanggal,
		Keterangan:  laporan.Keterangan,
		Pemasukan:   laporan.Pemasukan,
		Pengeluaran: laporan.Pengeluaran,
		Nota:        laporan.Nota,
		Saldo:       laporan.Saldo,
	}
}

func ConvertLaporanToListResponseDTO(laporan []model.LaporanKeuangan) []dto.LaporanKeuanganResponse {
	var laporanResponse []dto.LaporanKeuanganResponse
	for _, laporans := range laporan {
		laporanResponse = append(laporanResponse, ConvertLaporanToResponseDTO(laporans))
	}

	return laporanResponse
}
