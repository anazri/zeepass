package handlers

import (
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