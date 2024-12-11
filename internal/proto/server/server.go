package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"

	"thumbnail/internal/http/handlers"
	proto "thumbnail/internal/proto/proto"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedThumbnailServiceServer
	httpHandler *handlers.ThumbnailHandler
}

func NewServer(httpHandler *handlers.ThumbnailHandler) *Server {
	return &Server{httpHandler: httpHandler}
}

func (s *Server) GetThumbnail(ctx context.Context, req *proto.ThumbnailRequest) (*proto.ThumbnailResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	r := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: "http://localhost:8080/api/v1/thumbnail/"},
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: make(http.Header),
	}
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.httpHandler.GetThumbnail(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP handler returned error: %s", http.StatusText(resp.StatusCode))
	}

	var thumbnailResponse proto.ThumbnailResponse
	if err := json.NewDecoder(resp.Body).Decode(&thumbnailResponse); err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %w", err)
	}

	return &thumbnailResponse, nil
}

func (s *Server) GetThumbnails(ctx context.Context, req *proto.ThumbnailsRequest) (*proto.ThumbnailsResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	r := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: "http://localhost:8080/api/v1/thumbnail/"},
		Body:   io.NopCloser(bytes.NewReader(jsonData)),
		Header: make(http.Header),
	}

	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.httpHandler.GetThumbnail(w, r)
	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP handler returned error: %s", http.StatusText(resp.StatusCode))
	}

	var thumbnailsResponse proto.ThumbnailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&thumbnailsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode HTTP response: %w", err)
	}

	return &thumbnailsResponse, nil
}

func Run(address string, httpHandler *handlers.ThumbnailHandler) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	grpcServer := grpc.NewServer()
	server := NewServer(httpHandler)
	proto.RegisterThumbnailServiceServer(grpcServer, server)

	slog.Info("GRPC-server started successful")
	return grpcServer.Serve(lis)
}
