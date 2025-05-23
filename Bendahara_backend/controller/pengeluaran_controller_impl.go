package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type pengeluaranControllerImpl struct {
	PengeluaranService service.PengeluaranService
}

func NewPengeluaranController(pengeluaranService service.PengeluaranService) PengeluaranController {
	return &pengeluaranControllerImpl{
		PengeluaranService: pengeluaranService,
	}
}

func (p *pengeluaranControllerImpl) AddPengeluaran(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pengeluaranRequest := dto.PengeluaranRequest{}

	responseDTO, err := p.PengeluaranService.AddPengeluaran(r.Context(), r, pengeluaranRequest)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil menambahkan pengeluaran")
}

// UpdatePengeluaran implements PengeluaranController.
func (s *pengeluaranControllerImpl) UpdatePengeluaran(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	pengeluaranRequest := dto.PengeluaranRequest{}

	responseDTO, err := s.PengeluaranService.UpdatePengeluaran(r.Context(), r, pengeluaranRequest, id)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil memperbarui pengeluaran")
}

// GetPengeluaran implements PengeluaranController.
func (s *pengeluaranControllerImpl) GetPengeluaran(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Ambil parameter pagination dan date range dari query string
	page := util.StringToInt(r.URL.Query().Get("page"))
	pageSize := util.StringToInt(r.URL.Query().Get("page_size"))
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Set default jika tidak ada
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 15
	}

	var responseDTO dto.PengeluaranPaginationResponse
	var err error

	// Jika start_date dan end_date disediakan, gunakan GetPengeluaranByDateRange
	if startDate != "" && endDate != "" {
		responseDTO, err = s.PengeluaranService.GetPengeluaranByDateRange(r.Context(), startDate, endDate, page, pageSize)
	} else {
		// Jika tidak ada filter tanggal, gunakan GetPengeluaran seperti sebelumnya
		responseDTO, err = s.PengeluaranService.GetPengeluaran(r.Context(), page, pageSize)
	}

	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan pengeluaran dengan pagination")
}

// DeletePengeluaran implements PengeluaranController.
func (s *pengeluaranControllerImpl) DeletePengeluaran(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	responseDTO, err := s.PengeluaranService.DeletePengeluaran(r.Context(), id)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "pengeluaran yang berhasil dihapus")
}

// GetById implements PengeluaranController.
func (s *pengeluaranControllerImpl) GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	responseDTO, err := s.PengeluaranService.GetById(r.Context(), id)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan pengeluaran")
}
