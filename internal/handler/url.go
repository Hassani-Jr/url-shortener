package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Hassani-Jr/url-shortener/internal/service"
	"github.com/Hassani-Jr/url-shortener/internal/validator"
	"github.com/Hassani-Jr/url-shortener/pkg/logger/apperror"
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

type ShortenStats struct{
	LongURL string `json:"long_url"`
	ShortCode string `json:"short_code"`
	Timestamp time.Time `json:"timestamp"`
}

type DeleteRequest struct {
	Deleted bool `json:"deleted"`
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request){
	// parse request
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, apperror.BadRequest("Invalid JSON", err))
		return
	}

	validURL, err := validator.ValidateURL(req.URL)
	if err != nil{
		RespondError(w,err)
		return
	}

	// Call service layer
	shortCode, err := h.service.ShortenURL(r.Context(), validURL)
	if err != nil{
		RespondError(w,err)
		return
	}

	RespondJSON(w,http.StatusCreated,ShortenResponse{
		ShortCode: shortCode,
		ShortURL: "http://localhost:8080/" + shortCode,
	})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request){
	shortCode := r.PathValue("code")

	longUrl, err := h.service.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		RespondError(w, err)
		return
	}
	http.Redirect(w,r,longUrl,http.StatusMovedPermanently)
}

func (h *URLHandler) Stats(w http.ResponseWriter, r *http.Request){
	shortCode := r.PathValue("code")

	longUrl, err := h.service.GetOriginalURL(r.Context(),shortCode)
	if err != nil {
		RespondError(w, err)
		return
	}

	timestamp,err := h.service.GetTimeStamp(r.Context(),shortCode)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, ShortenStats{
		LongURL: longUrl,
		ShortCode: shortCode,
		Timestamp: timestamp,
	})
	 
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request){
	shortCode := r.PathValue("code")

	err := h.service.DeleteURL(r.Context(),shortCode)
	if err != nil {
		RespondError(w,err)
		return
	}

	RespondJSON(w, http.StatusNoContent,DeleteRequest{
		Deleted: true,
	})
}