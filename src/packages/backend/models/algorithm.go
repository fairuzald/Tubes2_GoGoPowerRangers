package models

// Struktur untuk request body
type AlgorithmRequest struct {
	Algorithm   string `json:"algorithm"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}
