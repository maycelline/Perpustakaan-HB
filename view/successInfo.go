package view

import (
	"encoding/json"
	"net/http"

	"../model"
)

func successProcess(w http.ResponseWriter) {
	var response model.SuccessResponse
	response.Status = 200
	response.Message = "Success"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode((response))
}
