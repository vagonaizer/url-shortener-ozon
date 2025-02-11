package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/vagonaizer/url-shortener-ozon/api/grpcserver"
	"github.com/vagonaizer/url-shortener-ozon/api/restserver"
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
	log.Printf("Запуск сервиса с BASE_URL: %s и REST-портом: %s", baseURL, restPort)

	// Определяем режим хранилища: "postgres" или иное (in-memory)
	storageMode := os.Getenv("STORAGE_MODE")
	var svc *service.URLShortenerService
	if storageMode == "postgres" {
		dbURL := os.Getenv("DATABASE_URL")
		pgStorage, err := storage.NewPostgresStorage(dbURL)
		if err != nil {
			log.Fatalf("Ошибка подключения к БД: %v", err)
		}
		svc = service.NewURLShortenerServiceWithStorage(pgStorage)
		log.Println("Используем PostgreSQL в качестве хранилища")
	} else {
		svc = service.NewURLShortenerService()
		log.Println("Используем in-memory хранилище")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Запуск gRPC сервера
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Ошибка при создании слушателя: %v", err)
		}
		grpcServer := grpc.NewServer()
		grpcserver.Register(grpcServer, svc)
		log.Println("gRPC сервер запущен на порту 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC сервер завершился с ошибкой: %v", err)
		}
	}()

	// Запуск REST сервера
	go func() {
		defer wg.Done()
		restHandler := restserver.New(svc)
		log.Printf("REST сервер запущен на порту %s", restPort)
		if err := http.ListenAndServe(":"+restPort, restHandler); err != nil {
			log.Fatalf("REST сервер завершился с ошибкой: %v", err)
		}
	}()

	wg.Wait()
}
