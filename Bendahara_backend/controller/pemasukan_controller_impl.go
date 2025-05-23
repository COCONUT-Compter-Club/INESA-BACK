package controller

import (
	// "encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type pemasukanControllerImpl struct {
	PemasukanService service.PemasukanService
}

func NewPemasukanController(pemasukanService service.PemasukanService) PemasukanController {
	return &pemasukanControllerImpl{
		PemasukanService: pemasukanService,
	}
}

// AddPemasukan implements PemasukanController.
func (s *pemasukanControllerImpl) AddPemasukan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pemasukanRequest := dto.PemasukanRequest{}
	// util.ReadFromRequestBody(r, &pemasukanRequest)

	responseDTO, err := s.PemasukanService.AddPemasukan(r.Context(), r, pemasukanRequest)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil ditambahkan")
}

// UpdatePemasukan implements PemasukanController.
func (s *pemasukanControllerImpl) UpdatePemasukan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	pemasukanRequest := dto.PemasukanRequest{}
	// util.ReadFromRequestBody(r, &pemasukanRequest)

	responseDTO, err := s.PemasukanService.UpdatePemasukan(r.Context(), r, pemasukanRequest, id)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil diperbarui")
}

// GetPemasukan implements PemasukanController.
func (s *pemasukanControllerImpl) GetPemasukan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Ambil parameter pagination dari query string
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

	var responseDTO dto.PemasukanPaginationResponse
	var err error

	// Jika start_date dan end_date disediakan, gunakan GetPemasukanByDateRange
	if startDate != "" && endDate != "" {
		responseDTO, err = s.PemasukanService.GetPemasukanByDateRange(r.Context(), startDate, endDate, page, pageSize)
	} else {
		// Jika tidak ada filter tanggal, gunakan GetPemasukan seperti sebelumnya
		responseDTO, err = s.PemasukanService.GetPemasukan(r.Context(), page, pageSize)
	}

	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan pemasukan dengan pagination")
}

// DeletePemasukan implements PemasukanController.
func (s *pemasukanControllerImpl) DeletePemasukan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	responseDTO, err := s.PemasukanService.DeletePemasukan(r.Context(), id)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil menghapus")
}

// GetById implements PemasukanController.
func (s *pemasukanControllerImpl) GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	responseDTO, err := s.PemasukanService.GetById(r.Context(), id)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan")
}