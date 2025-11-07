package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/Hassani-Jr/url-shortener/internal/storage"
	"github.com/Hassani-Jr/url-shortener/pkg/logger/apperror"
)

type ShortenerService struct{
	storage storage.URLStorage
}

func NewShortenerService(storage storage.URLStorage) *ShortenerService{
	return &ShortenerService{storage: storage}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, longURL string) (string, error) {
	// check if context is already cancelled
	if ctx.Err() != nil {
		return "",apperror.Internal("Request cancelled", ctx.Err())
	}

	// Generate short code
	shortCode, err := generateShortCode()
	if err != nil {
		return "", apperror.Internal("Failed to generate code", err)
	}
	
	

	// Save to storage
	if err := s.storage.Save(ctx, &storage.URL{ShortCode: shortCode,LongURL: longURL,CreatedAt: time.Now()}); err != nil {
		return "", err
	}

	return shortCode, nil
}

func (s *ShortenerService) GetOriginalURL(ctx context.Context, shortCode string) (string,error){
	// Context timeout example
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()

	url, err := s.storage.Get(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if url.LongURL == ""{
		return "", apperror.NotFound("Short URL not found")
	}

	return url.LongURL, nil
}

func (s *ShortenerService) GetTimeStamp(ctx context.Context, shortCode string) (time.Time, error){
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()

	url, err := s.storage.Get(ctx,shortCode)
	if err != nil {
		return time.Time{}, err
	}

	if url.CreatedAt.IsZero() {
		return time.Time{}, apperror.NotFound("Timestamp not found")
	}

	return url.CreatedAt,nil
}

func (s *ShortenerService) DeleteURL(ctx context.Context, shortCode string)(error) {
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()

	err := s.storage.Delete(ctx,shortCode)
	if err != nil{
		return err
	}

	return nil
}

func generateShortCode() (string, error){
	b := make([]byte,6)
	if _,err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:8], nil
}