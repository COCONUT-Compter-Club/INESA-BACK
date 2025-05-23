package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/dto"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
	"github.com/syrlramadhan/api-bendahara-inovdes/util"
)

type adminControllerImpl struct {
	AdminService service.AdminService
}

func NewAdminController(adminService service.AdminService) AdminController {
	return adminControllerImpl{
		AdminService: adminService,
	}
}

// SignUp implements AdminController.
func (a adminControllerImpl) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	adminRequest := dto.AdminRequest{}
	util.ReadFromRequestBody(r, &adminRequest)

	responseDTO, err := a.AdminService.SignUp(r.Context(), adminRequest)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "pendaftaran berhasil")
}

// SignIn implements AdminController.
func (a adminControllerImpl) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginRequest := dto.LoginRequest{}
	util.ReadFromRequestBody(r, &loginRequest)

	responseDTO, err := a.AdminService.SignIn(r.Context(), loginRequest)
	if err != nil {
		helper.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil masuk")
}

// FindByNik implements AdminController.
func (a adminControllerImpl) FindByNik(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	nik := ps.ByName("nik")
	responseDTO, err := a.AdminService.GetAdminByNik(r.Context(), nik)
	if err != nil {
		helper.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "dapatkan admin berhasil")
}