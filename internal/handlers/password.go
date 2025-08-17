package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/anazri/zeepass/internal/models"
	"github.com/anazri/zeepass/internal/services"
)

func PasswordGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/password-generator.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := models.PageData{
		Title: "Password Generator - ZeePass",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func GeneratePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var opts services.PasswordOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		// Try to parse as form data for backward compatibility
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing request", http.StatusBadRequest)
			return
		}

		// Parse form values
		lengthStr := r.FormValue("length")
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			length = 12 // default
		}

		opts = services.PasswordOptions{
			Length:       length,
			UseNumbers:   r.FormValue("use_numbers") == "true",
			UseUppercase: r.FormValue("use_uppercase") == "true",
			UseLowercase: r.FormValue("use_lowercase") == "true",
			UseSymbols:   r.FormValue("use_symbols") == "true",
			Type:         r.FormValue("type"),
		}

		// Default to numbers if nothing selected
		if !opts.UseNumbers && !opts.UseUppercase && !opts.UseLowercase && !opts.UseSymbols {
			opts.UseNumbers = true
		}
	}

	// Generate password
	password, err := services.GeneratePassword(opts)
	if err != nil {
		http.Error(w, "Error generating password", http.StatusInternalServerError)
		log.Printf("Password generation error: %v", err)
		return
	}

	// Calculate strength
	strength := services.CalculatePasswordStrength(password)

	// Return JSON response
	response := map[string]interface{}{
		"password": password,
		"strength": strength,
		"length":   len(password),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}