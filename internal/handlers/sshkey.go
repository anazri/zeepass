package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/anazri/zeepass/internal/models"
	"github.com/anazri/zeepass/internal/services"
)

func SSHKeyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ssh-key.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := models.PageData{
		Title: "SSH Key - ZeePass",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func GenerateSSHKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var opts services.SSHKeyOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		http.Error(w, "Error parsing request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate options
	if err := services.ValidateSSHKeyOptions(opts); err != nil {
		http.Error(w, "Invalid options: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Generate SSH key pair
	keyPair, err := services.GenerateSSHKey(opts)
	if err != nil {
		log.Printf("SSH key generation error: %v", err)
		http.Error(w, "Error generating SSH key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyPair)
}