package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Hassani-Jr/url-shortener/internal/service"
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
	//Validate
	if req.URL == ""{
		RespondError(w, apperror.BadRequest("URL is required", nil))
		return 
	}

	if len(req.URL) > 2000 {
		RespondError(w, apperror.BadRequest("URL is too long",nil))
		return
	}

	if !strings.EqualFold(req.URL[0:8],"https://") && !strings.EqualFold(req.URL[0:7],"http://") {
		RespondError(w, apperror.BadRequest("Unsafe URL", nil))
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

	RespondJSON(w, http.StatusAccepted, ShortenStats{
		LongURL: longUrl,
		ShortCode: shortCode,
		Timestamp: timestamp,
	})
	 
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request){
	shortCode := r.PathValue("code")

	deleted ,err := h.service.DeleteURL(r.Context(),shortCode)
	if err != nil {
		RespondError(w,err)
		return
	}

	RespondJSON(w, http.StatusOK,DeleteRequest{
		Deleted: deleted,
	})
}