package storage

import (
	"context"
	"sync"
	"time"

	"github.com/Hassani-Jr/url-shortener/pkg/logger/apperror"
)

type URLStorage interface {
	Save(ctx context.Context, url *URL) error
	Get(ctx context.Context, shortCode string) (*URL, error)
	Delete(ctx context.Context,shortCode string) (error)
}

type URL struct{
	ShortCode string
	LongURL string
	CreatedAt time.Time
}

// MemoryStorage is an in-memory implementation (testing/dev)
type MemoryStorage struct {
	mu sync.RWMutex
	urls map[string]URL
}

func NewMemoryStorage() *MemoryStorage{
	return &MemoryStorage{
		urls : make(map[string]URL),
	}
}
 func (m *MemoryStorage) Save(ctx context.Context, url *URL) error {
	if ctx.Err() != nil{
		return apperror.Internal("Context Cancelled", ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Duplicates
	if _,exists := m.urls[url.ShortCode]; exists {
		return apperror.BadRequest("Short code already exists", nil)
	}
	 m.urls[url.ShortCode] = URL{url.ShortCode,url.LongURL,url.CreatedAt}

	 return nil
 }

 func (m *MemoryStorage) Get(ctx context.Context, shortCode string) (*URL, error) {
	if ctx.Err() != nil {
		return &URL{}, apperror.Internal("Context cancelled", ctx.Err())
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	url,exists := m.urls[shortCode]
	if !exists {
		return &URL{},apperror.NotFound("Short code not found")
	}

	return  &url,nil
 }

 func (m *MemoryStorage) Delete(ctx context.Context, shortCode string) (error) {
	if ctx.Err() != nil{
		return apperror.Internal("Context Cancelled", ctx.Err())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	_,exists := m.urls[shortCode]
	if !exists {
		return apperror.NotFound("Short code not found")
	}
	delete(m.urls,shortCode)

	return nil
 }