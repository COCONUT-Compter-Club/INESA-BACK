package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
	"github.com/syrlramadhan/api-bendahara-inovdes/service"
)

type transactionControllerImpl struct {
	TransactionService service.TransactionService
}

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionControllerImpl{
		TransactionService: transactionService,
	}
}

// GetAllTransaction implements TransactionController.
func (t *transactionControllerImpl) GetAllTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseDTO, err := t.TransactionService.GetAllTransaction(r.Context())
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan semua transaksi")
}

// GetLastTransaction implements TransactionController.
func (t *transactionControllerImpl) GetLastTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseDTO, err := t.TransactionService.GetLastTransaction(r.Context())
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSONSuccess(w, responseDTO, "berhasil mendapatkan transaksi")
}
