package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/anazri/zeepass/internal/models"
)

// createPageData creates a PageData struct with common fields populated
func createPageData(title string) models.PageData {
	return models.PageData{
		Title:                   title,
		RecaptchaSiteKey:        os.Getenv("RECAPTCHA_SITE_KEY"),
		CloudflareAnalyticsToken: os.Getenv("CLOUDFLARE_ANALYTICS_TOKEN"),
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := createPageData("ZeePass - Encrypt your data easily")

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func TextEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/text-encryption.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := createPageData("Text Encryption - ZeePass")

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func FileEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/file-encryption.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := createPageData("File Encryption - ZeePass")

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))).ServeHTTP(w, r)
}
