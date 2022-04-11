package controllers

import (
	"Perpustakaan-HB/model"
	"encoding/json"
	"net/http"
)

func sendSuccessResponse(w http.ResponseWriter, message string, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var response model.SuccessResponse
	response.Status = http.StatusOK
	response.Message = message
	response.Data = value
	json.NewEncoder(w).Encode(response)
}

func sendBadRequestResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	var response model.ErrorResponse
	response.Status = http.StatusBadRequest
	response.Message = errorMessage
	json.NewEncoder(w).Encode(response)
}

func sendUnauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	var response model.ErrorResponse
	response.Status = http.StatusUnauthorized
	response.Message = "Unauthorized Access"
	json.NewEncoder(w).Encode(response)
}

func sendNotFoundResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	var response model.ErrorResponse
	response.Status = http.StatusNotFound
	response.Message = errorMessage
	json.NewEncoder(w).Encode(response)
}

func sendServerErrorResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	var response model.ErrorResponse
	response.Status = http.StatusInternalServerError
	response.Message = errorMessage
	json.NewEncoder(w).Encode(response)
}

func sendSuccessResponseWithoutData(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	var response model.SuccessResponse
	response.Status = http.StatusOK
	response.Message = message
	json.NewEncoder(w).Encode(response)
}
