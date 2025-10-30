package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Hassani-Jr/url-shortener/pkg/logger/apperror"
)

type ErrorResponse struct{
	Error string `json:"error"`
	Code string `json:"code"`
	Message string `json:"message"`
}

// RespondError sends a JSON error response
func RespondError(w http.ResponseWriter, err error){

	//Check to see if its a custom err
	if appErr, ok := err.(*apperror.AppError); ok {
		//Log internal details
		if appErr.StatusCode >= 500{
			log.Printf("Internal error: %v", appErr.Err)
		}
	

	RespondJSON(w, appErr.StatusCode, ErrorResponse{
		Error: appErr.Message,
		Code: appErr.Code,
		Message: appErr.Message,
	})
	return
	}
	log.Printf("Unexpected error: %v", err)
	RespondJSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: "Internal server error",
		Code: "INTERNAL_ERROR",
		Message: "Something went wrong",
	})
}

// sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}