package storage

import (
	"context"
	"sync"
	"time"

	"github.com/Hassani-Jr/url-shortener/pkg/logger/apperror"
)

type URLStorage interface {
	Save(ctx context.Context, shortCode, longURL string) error
	Get(ctx context.Context, shortCode string) (UrlMap, error)
	Delete(ctx context.Context,shortCode string) (bool,error)
}

type UrlMap struct{
	Value string
	Time time.Time
}

// MemoryStorage is an in-memory implementation (testing/dev)
type MemoryStorage struct {
	mu sync.RWMutex
	urls map[string]UrlMap
}

func NewMemoryStorage() *MemoryStorage{
	return &MemoryStorage{
		urls : make(map[string]UrlMap),
	}
}
 func (m *MemoryStorage) Save(ctx context.Context, shortCode, longURL string) error {
	if ctx.Err() != nil{
		return apperror.Internal("Context Cancelled", ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Duplicates
	if _,exists := m.urls[shortCode]; exists {
		return apperror.BadRequest("Short code already exists", nil)
	}
	 m.urls[shortCode] = UrlMap{longURL,time.Now()}

	 return nil
 }

 func (m *MemoryStorage) Get(ctx context.Context, shortCode string) (UrlMap, error) {
	if ctx.Err() != nil {
		return UrlMap{}, apperror.Internal("Context cancelled", ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	url,exists := m.urls[shortCode]
	if !exists {
		return UrlMap{},apperror.NotFound("Short code not found")
	}

	return  url,nil
 }

 func (m *MemoryStorage) Delete(ctx context.Context, shortCode string) (bool,error) {
	if ctx.Err() != nil{
		return false ,apperror.Internal("Context Cancelled", ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	_,exists := m.urls[shortCode]
	if !exists {
		return false,apperror.NotFound("Short code not found")
	}
	delete(m.urls,shortCode)

	return true,nil
 }