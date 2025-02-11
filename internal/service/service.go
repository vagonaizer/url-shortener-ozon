package service

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sync"
	"time"

	"github.com/vagonaizer/url-shortener-ozon/internal/storage"
)

// Storage описывает операции сохранения и получения URL.
type Storage interface {
	storage.Storage
}

// InMemoryStorage – реализация Storage на базе карты.
type InMemoryStorage struct {
	urlMap map[string]string
	mu     sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urlMap: make(map[string]string),
	}
}

func (s *InMemoryStorage) Save(shortURL, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urlMap[shortURL] = originalURL
	return nil
}

func (s *InMemoryStorage) Get(shortURL string) (string, bool) {
	s.mu.RLock()
	originalURL, ok := s.urlMap[shortURL]
	s.mu.RUnlock()
	return originalURL, ok
}

// URLShortenerService использует хранилище для работы.
type URLShortenerService struct {
	storage Storage
}

// Конструкторы
func NewURLShortenerService() *URLShortenerService {
	return &URLShortenerService{
		storage: NewInMemoryStorage(),
	}
}

func NewURLShortenerServiceWithStorage(store Storage) *URLShortenerService {
	return &URLShortenerService{
		storage: store,
	}
}

// ShortenURL генерирует сокращённую ссылку длиной 10 символов.
func (s *URLShortenerService) ShortenURL(originalURL string) string {
	// Если уже есть сокращённая ссылка для данного URL, возвращаем её.
	// (В данной простейшей реализации поиск по значению не реализован.
	//  Можно сделать дополнительную логику для обхода дублирования.)

	// Генерация на основе хеша + случайная соль
	h := sha256.New()
	salt := generateRandomString(5)
	h.Write([]byte(originalURL + salt))
	candidate := hex.EncodeToString(h.Sum(nil))[:10] // берем первые 10 символов

	// Сохраняем, игнорируя ошибку, если уже есть такой ключ.
	_ = s.storage.Save(candidate, originalURL)
	return candidate
}

func (s *URLShortenerService) GetOriginalURL(shortURL string) (string, bool) {
	return s.storage.Get(shortURL)
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
