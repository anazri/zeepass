package handlers

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/anazri/zeepass/internal/models"
)

func Base64Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/base64.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := models.PageData{
		Title: "Base64 - ZeePass",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func Base64EncodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	dataType := r.FormValue("type")
	var result string

	if dataType == "file" {
		// Handle file encoding
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read file content
		fileData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file content", http.StatusInternalServerError)
			return
		}

		// Encode to base64
		result = base64.StdEncoding.EncodeToString(fileData)
	} else {
		// Handle text encoding
		text := r.FormValue("text")
		if text == "" {
			http.Error(w, "No text provided", http.StatusBadRequest)
			return
		}

		// Encode to base64
		result = base64.StdEncoding.EncodeToString([]byte(text))
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"result":  result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Base64DecodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	base64Data := r.FormValue("data")
	if base64Data == "" {
		http.Error(w, "No base64 data provided", http.StatusBadRequest)
		return
	}

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		http.Error(w, "Invalid base64 data", http.StatusBadRequest)
		return
	}

	dataType := r.FormValue("type")

	if dataType == "file" {
		// Return as file download
		filename := r.FormValue("filename")
		if filename == "" {
			filename = "decoded_file"
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		w.Header().Set("Content-Length", string(rune(len(decoded))))
		w.Write(decoded)
	} else {
		// Return as text in JSON response
		response := map[string]interface{}{
			"success": true,
			"result":  string(decoded),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}