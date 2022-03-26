package controllers

import (
	"Tools/model"
	"encoding/json"
	"net/http"
)

// func sendSuccessResponse(w http.ResponseWriter, message string, books []Book) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var response BorrowData
// 	response.Status = 200
// 	response.Message = message
// 	response.Data = users
// 	json.NewEncoder(w).Encode(response)
// }

func sendSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	var response model.SuccessResponse
	response.Status = 200
	response.Message = message
	json.NewEncoder(w).Encode(response)
}

func sendBadRequestResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	var response model.ErrorResponse
	response.Status = 400
	response.Message = errorMessage
	json.NewEncoder(w).Encode(response)
}

// func sendUnauthorizedResponse(w http.ResponseWriter) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusUnauthorized)
// 	var response model.ErrorResponse
// 	response.Status = 401
// 	response.Message = "Unauthorized Access"
// 	json.NewEncoder(w).Encode(response)
// }

func sendNotFoundResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	var response model.ErrorResponse
	response.Status = 404
	response.Message = errorMessage
	json.NewEncoder(w).Encode(response)
}
