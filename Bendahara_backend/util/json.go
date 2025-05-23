package util

import (
	"encoding/json"
	"net/http"

	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
)

func ReadFromRequestBody(request *http.Request, result interface{}) (error) {
	var writer http.ResponseWriter
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		helper.WriteJSONError(writer, http.StatusOK, "gagal mengkodekan respons JSON")
	}

	return err
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		helper.WriteJSONError(writer, http.StatusOK, "gagal mengkodekan respons JSON")
	}
}