package main

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vagonaizer/url-shortener-ozon/api/grpcserver"
	"github.com/vagonaizer/url-shortener-ozon/api/restserver"
	"github.com/vagonaizer/url-shortener-ozon/internal/logger"
	"github.com/vagonaizer/url-shortener-ozon/internal/service"
	"github.com/vagonaizer/url-shortener-ozon/internal/storage"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	restPort := os.Getenv("PORT")
	if restPort == "" {
		restPort = "8080"
	}

	baseURL := os.Getenv("BASE_URL")
	logger.InfoLogger.Printf("Запуск сервиса с BASE_URL: %s и REST-портом: %s", baseURL, restPort)

	// Определяем режим хранилища: "postgres" или иное (in-memory)
	storageMode := os.Getenv("STORAGE_MODE")
	var svc *service.URLShortenerService
	if storageMode == "postgres" {
		dbURL := os.Getenv("DATABASE_URL")
		pgStorage, err := storage.NewPostgresStorage(dbURL)
		if err != nil {
			logger.ErrorLogger.Fatalf("Ошибка подключения к БД: %v", err)
		}
		svc = service.NewURLShortenerServiceWithStorage(pgStorage)
		logger.InfoLogger.Println("Используем PostgreSQL в качестве хранилища")
	} else {
		svc = service.NewURLShortenerService()
		logger.InfoLogger.Println("Используем in-memory хранилище")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Запуск gRPC сервера в основной горутине
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.ErrorLogger.Fatalf("Ошибка при создании слушателя: %v", err)
	}
	grpcServer := grpc.NewServer()
	grpcserver.Register(grpcServer, svc)
	logger.InfoLogger.Println("gRPC сервер запущен на порту 50051")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.ErrorLogger.Fatalf("gRPC сервер завершился с ошибкой: %v", err)
		}
	}()

	// Запуск REST сервера в основной горутине
	restHandler := restserver.New(svc)
	logger.InfoLogger.Printf("REST сервер запущен на порту %s", restPort)
	if err := http.ListenAndServe(":"+restPort, restHandler); err != nil {
		logger.ErrorLogger.Fatalf("REST сервер завершился с ошибкой: %v", err)
	}

	wg.Wait()
}
