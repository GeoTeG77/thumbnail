package models

type ThumbnailResponse struct {
	URL     string `json:"url"`
	Preview []byte `json:"preview"`
}
type ThumbnailRequest struct {
	URL string `json:"url"`
}
