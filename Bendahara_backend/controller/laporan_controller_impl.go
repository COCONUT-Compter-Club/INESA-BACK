package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type laporanKeuanganControllerImpl struct {
	LaporanService service.LaporanKeuanganService
}

func NewLaporanKeuanganController(laporanService service.LaporanKeuanganService) LaporanKeuanganController {
	return &laporanKeuanganControllerImpl{
		LaporanService: laporanService,
	}
}

// GetAllLaporan implements LaporanKeuanganController.
func (l *laporanKeuanganControllerImpl) GetAllLaporan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseDTO, err := l.LaporanService.GetAllLaporan(r.Context())
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan semua data")
}

// GetLastBalance implements LaporanKeuanganController.
func (l *laporanKeuanganControllerImpl) GetLastBalance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	saldo, err := l.LaporanService.GetLastBalance(r.Context())
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := dto.ListResponseSaldo{
		Code: http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Saldo: saldo,
		Message: "berhasil mendapatkan saldo terakhir",
	}
	
	util.WriteToResponseBody(w, response)
}

// GetTotalExpenditure implements LaporanKeuanganController.
func (l *laporanKeuanganControllerImpl) GetTotalExpenditure(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	totalPengeluaran, err := l.LaporanService.GetTotalExpenditure(r.Context())
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ListResponseSaldo{
		Code: http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Saldo: int64(totalPengeluaran),
		Message: "berhasil mendapatkan total pengeluaran",
	}

	util.WriteToResponseBody(w, response)
}

// GetTotalIncome implements LaporanKeuanganController.
func (l *laporanKeuanganControllerImpl) GetTotalIncome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	totalPemasukan, err := l.LaporanService.GetTotalIncome(r.Context())
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ListResponseSaldo{
		Code: http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Saldo: int64(totalPemasukan),
		Message: "berhasil mendapatkan total pemasukan",
	}

	util.WriteToResponseBody(w, response)
}

func (l *laporanKeuanganControllerImpl) GetLaporanByDateRange(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Ambil parameter query dari URL
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	// Validasi parameter
	if startDate == "" || endDate == "" {
		helper.WriteJSONError(w, http.StatusBadRequest, "tanggal mulai dan tanggal akhir diperlukan")
		return
	}

	// Panggil service untuk mendapatkan data laporan berdasarkan rentang tanggal
	responseDTO, err := l.LaporanService.GetLaporanByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Kirim respons JSON
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan laporan berdasarkan rentang tanggal")
}