package grpcserver

import (
	"context"

	pb "github.com/vagonaizer/url-shortener-ozon/api/proto"
	"github.com/vagonaizer/url-shortener-ozon/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server реализует pb.URLShortenerServer.
type Server struct {
	pb.UnimplementedURLShortenerServer
	svc *service.URLShortenerService
}

// Register регистрирует сервер gRPC.
func Register(srv *grpc.Server, svc *service.URLShortenerService) {
	pb.RegisterURLShortenerServer(srv, &Server{svc: svc})
}

// ShortenURL создает короткую ссылку для переданного URL.
func (s *Server) ShortenURL(ctx context.Context, req *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	shortURL := s.svc.ShortenURL(req.OriginalUrl)
	return &pb.ShortenURLResponse{ShortUrl: shortURL}, nil
}

// GetOriginalURL возвращает исходный URL по короткой ссылке.
func (s *Server) GetOriginalURL(ctx context.Context, req *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	originalURL, ok := s.svc.GetOriginalURL(req.ShortUrl)
	if !ok {
		return nil, status.Error(codes.NotFound, "URL не найден")
	}
	return &pb.GetOriginalURLResponse{OriginalUrl: originalURL}, nil
}
