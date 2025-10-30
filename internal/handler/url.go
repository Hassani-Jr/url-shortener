package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Hassani-Jr/url-shortener/pkg/logger/apperror"
	"github.com/Hassani-Jr/url-shortener/internal/service"
)

type URLHandler struct{
	service *service.ShortenerService
}

func NewURLHandler(service *service.ShortenerService) *URLHandler {
	return &URLHandler{service: service}
}

type ShortenRequest struct {
	URL string`json:"url"`
}

type ShortenResponse struct{
	ShortCode string `json:"short_code"`
	ShortURL string `json:"short_url"`
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request){
	// parse request
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, apperror.BadRequest("Invalid JSON", err))
		return
	}

	//Validate
	if req.URL == ""{
		RespondError(w, apperror.BadRequest("URL is required", nil))
		return 
	}

	// Call service layer
	shortCode, err := h.service.ShortenURL(r.Context(), req.URL)
	if err != nil{
		RespondError(w,err)
		return
	}

	RespondJSON(w,http.StatusCreated,ShortenResponse{
		ShortCode: shortCode,
		ShortURL: "http://localhost:8080/" + shortCode,
	})
}