package router

import (
	"net/http"
	
	"thumbnail/internal/http/handlers"
)

func NewRouter(handler *handlers.ThumbnailHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/thumbnail/", handler.GetThumbnail)

	return mux
}
