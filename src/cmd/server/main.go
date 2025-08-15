package main

import (
	"log"
	"net/http"

	"github.com/anazri/zeepass/internal/handlers"
	"github.com/anazri/zeepass/internal/services"
)

func main() {
	services.InitRedis()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/text-encryption", handlers.TextEncryptionHandler)
	http.HandleFunc("/encrypt-text", handlers.EncryptTextHandler)
	http.HandleFunc("/file-encryption", handlers.FileEncryptionHandler)
	http.HandleFunc("/encrypt-file", handlers.EncryptFileHandler)
	http.HandleFunc("/chat-encryption", handlers.ChatEncryptionHandler)
	http.HandleFunc("/ws/chat", handlers.ChatWebSocketHandler)
	http.HandleFunc("/password-generator", handlers.PasswordGeneratorHandler)
	http.HandleFunc("/generate-password", handlers.GeneratePasswordHandler)
	http.HandleFunc("/base64", handlers.Base64Handler)
	http.HandleFunc("/base64-encode", handlers.Base64EncodeHandler)
	http.HandleFunc("/base64-decode", handlers.Base64DecodeHandler)
	http.HandleFunc("/ssh-key", handlers.SSHKeyHandler)
	http.HandleFunc("/generate-ssh-key", handlers.GenerateSSHKeyHandler)
	http.HandleFunc("/view/", handlers.ViewEncryptedHandler)
	http.HandleFunc("/view-file/", handlers.ViewEncryptedFileHandler)
	http.HandleFunc("/contact", handlers.HandleContact)
	http.HandleFunc("/static/", handlers.StaticHandler)

	log.Println("ZeePass server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
