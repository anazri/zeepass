package models

import "time"

type PageData struct {
	Title string
}

type EncryptedData struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	PIN       string     `json:"pin,omitempty"`
	Lifetime  string     `json:"lifetime"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	ViewCount int        `json:"view_count"`
	MaxViews  int        `json:"max_views"`
}

type EncryptionRequest struct {
	Text     string `json:"text"`
	PIN      string `json:"pin"`
	Lifetime string `json:"lifetime"`
}

type EncryptedFileData struct {
	ID         string     `json:"id"`
	Content    []byte     `json:"content"`
	FileName   string     `json:"file_name"`
	FileSize   int64      `json:"file_size"`
	MimeType   string     `json:"mime_type"`
	PIN        string     `json:"pin,omitempty"`
	Lifetime   string     `json:"lifetime"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	ViewCount  int        `json:"view_count"`
	MaxViews   int        `json:"max_views"`
}

type EncryptionResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	URL     string `json:"url"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}