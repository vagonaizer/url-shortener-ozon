package service_test

import (
	"sync"
	"testing"

	"github.com/vagonaizer/url-shortener-ozon/internal/service"
)

func TestShortenURLAndGetOriginalURL(t *testing.T) {
	svc := service.NewURLShortenerService()
	originalURL := "https://example.com"

	// Проверяем, что метод ShortenURL возвращает не пустую строку.
	shortURL := svc.ShortenURL(originalURL)
	if shortURL == "" {
		t.Fatalf("ShortenURL вернул пустую строку")
	}

	// Для одного и того же URL должна формироваться одна и та же сокращенная ссылка.
	shortURL2 := svc.ShortenURL(originalURL)
	if shortURL != shortURL2 {
		t.Errorf("Ожидается одинаковая сокращенная ссылка для одного URL, получены %s и %s", shortURL, shortURL2)
	}

	// Проверяем, что по сокращенной ссылке корректно возвращается оригинальный URL.
	retOriginal, exists := svc.GetOriginalURL(shortURL)
	if !exists {
		t.Fatalf("GetOriginalURL не нашёл оригинальный URL для: %s", shortURL)
	}
	if retOriginal != originalURL {
		t.Errorf("Ожидался %s, получено %s", originalURL, retOriginal)
	}
}

func TestConcurrentAccess(t *testing.T) {
	const numGoroutines = 100
	svc := service.NewURLShortenerService()
	originalURL := "https://concurrent.example.com"
	var wg sync.WaitGroup

	// Запускаем большое количество горутин для вызова ShortenURL.
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc.ShortenURL(originalURL)
		}()
	}
	wg.Wait()

	// Проверяем, что метод GetOriginalURL возвращает правильный оригинальный URL.
	shortURL := svc.ShortenURL(originalURL)
	retURL, exists := svc.GetOriginalURL(shortURL)
	if !exists {
		t.Fatalf("После параллельного доступа GetOriginalURL не нашёл URL для: %s", shortURL)
	}
	if retURL != originalURL {
		t.Errorf("Ожидается %s, получено %s", originalURL, retURL)
	}
}
