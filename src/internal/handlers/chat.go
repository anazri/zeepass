package handlers

import (
	"html/template"
	"net/http"
	"log"

	"github.com/anazri/zeepass/internal/services"
)

// ChatEncryptionHandler handles the chat encryption page
func ChatEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/chat-encryption.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Chat Encryption - ZeePass",
	}

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