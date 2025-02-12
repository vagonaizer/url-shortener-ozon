package storage

import (
	"sync"

	"github.com/vagonaizer/url-shortener-ozon/internal/storage/models"
)

// InMemoryStorage - структура для хранения данных в памяти
type InMemoryStorage struct {
	urlMap map[string]models.URLData
	mu     sync.RWMutex
}

// NewInMemoryStorage - конструктор для создания нового in-memory хранилища
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urlMap: make(map[string]models.URLData),
	}
}

// Save - метод для сохранения данных
func (s *InMemoryStorage) Save(shortCode string, data models.URLData) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urlMap[shortCode] = data
	return nil
}

// Get - метод для получения данных по короткому коду
func (s *InMemoryStorage) Get(shortCode string) (models.URLData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, exists := s.urlMap[shortCode]
	return data, exists
}
