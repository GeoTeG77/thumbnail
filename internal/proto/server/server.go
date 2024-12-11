package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"google.golang.org/grpc"

	"thumbnail/internal/http/handlers"
	"thumbnail/internal/models"
	"thumbnail/internal/proto/proto"
)

type Server struct {
	proto.UnimplementedThumbnailServiceServer
	httpHandler *handlers.ThumbnailHandler
}

func NewServer(httpHandler *handlers.ThumbnailHandler) *Server {
	return &Server{httpHandler: httpHandler}
}

func ConvertProtoToModelRequest(protoReq *proto.ThumbnailRequest) *models.ThumbnailRequest {
	return &models.ThumbnailRequest{
		URL: protoReq.GetUrl(),
	}
}

func ConvertModelToProtoRequest(modelReq *models.ThumbnailRequest) *proto.ThumbnailRequest {
	return &proto.ThumbnailRequest{
		Url: modelReq.URL,
	}
}

func ConvertModelToProtoResponse(modelResp *models.ThumbnailResponse) *proto.ThumbnailResponse {
	return &proto.ThumbnailResponse{
		Url:       modelResp.URL,
		ImageData: modelResp.Preview,
	}
}

func (s *Server) GetThumbnail(ctx context.Context, req *proto.ThumbnailRequest) (*proto.ThumbnailResponse, error) {
	httpRequest := ConvertProtoToModelRequest(req)
	jsonData, err := json.Marshal([]*models.ThumbnailRequest{httpRequest})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	r := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: os.Getenv("ENDPOINT_ADDRES")},
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

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	if len(imgData) == 0 {
		return nil, fmt.Errorf("received empty image data")
	}

	return &proto.ThumbnailResponse{
		ImageData: imgData,
		Url:       req.GetUrl(),
	}, nil
}

func (s *Server) GetThumbnails(ctx context.Context, req *proto.ThumbnailsRequest) (*proto.ThumbnailsResponse, error) {
	var httpRequests []*models.ThumbnailRequest
	for _, url := range req.GetUrls() {
		httpRequests = append(httpRequests, &models.ThumbnailRequest{URL: url})
	}

	jsonData, err := json.Marshal(httpRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	r := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: os.Getenv("ENDPOINT_ADDRES")},
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
	boundary := "myBoundary"
	reader := multipart.NewReader(resp.Body, boundary)
	var protoResponse []*proto.ThumbnailResult

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read multipart part: %w", err)
		}

		if part.Header.Get("Content-Type") == "image/jpeg" {
			imgData, err := io.ReadAll(part)
			if err != nil {
				return nil, fmt.Errorf("failed to read image data: %w", err)
			}

			protoResponse = append(protoResponse, &proto.ThumbnailResult{
				Url:       part.FileName(),
				ImageData: imgData,
			})
		}
	}

	return &proto.ThumbnailsResponse{Results: protoResponse}, nil
}

func Run(address string, httpHandler *handlers.ThumbnailHandler) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	grpcServer := grpc.NewServer()
	server := NewServer(httpHandler)
	proto.RegisterThumbnailServiceServer(grpcServer, server)

	slog.Info("GRPC-server started successfully")
	return grpcServer.Serve(lis)
}
