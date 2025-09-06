package handlers

import (
	"html/template"
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

func SurveyHandler(w http.ResponseWriter, r *http.Request) {
	data := createPageData("Survey - ZeePass Feedback")

	tmpl, err := template.ParseFiles("templates/survey.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
