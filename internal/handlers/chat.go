package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/anazri/zeepass/internal/models"
	"github.com/anazri/zeepass/internal/services"
)

// createPageData creates a PageData struct with common fields populated
func createPageData(title string) models.PageData {
	return models.PageData{
		Title:                   title,
		RecaptchaSiteKey:        os.Getenv("RECAPTCHA_SITE_KEY"),
		CloudflareAnalyticsToken: os.Getenv("CLOUDFLARE_ANALYTICS_TOKEN"),
	}
}

// ChatEncryptionHandler handles the chat encryption page
func ChatEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/chat-encryption.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := createPageData("Chat Encryption - ZeePass")

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

// ChatWebSocketHandler handles WebSocket connections for real-time chat
func ChatWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	chatService := services.GetChatService()
	chatService.HandleWebSocket(w, r)
}
