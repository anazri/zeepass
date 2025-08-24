package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type SurveyResponse struct {
	ID                 string    `json:"id"`
	Timestamp          time.Time `json:"timestamp"`
	Likelihood         string    `json:"likelihood"`
	Tools              []string  `json:"tools"`
	UseCase            string    `json:"use_case"`
	BusinessSector     string    `json:"business_sector"`
	EnterpriseInterest string    `json:"enterprise_interest"`
	Concerns           string    `json:"concerns"`
	FeatureRequest     string    `json:"feature_request"`
	NPS                int       `json:"nps"`
	Email              string    `json:"email,omitempty"`
	Name               string    `json:"name,omitempty"`
	Updates            bool      `json:"updates"`
	IPAddress          string    `json:"ip_address"`
}

func HandleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Generate unique ID
	id := fmt.Sprintf("survey_%d_%d", time.Now().Unix(), time.Now().Nanosecond())

	// Parse NPS score
	npsScore := 0
	if npsStr := strings.TrimSpace(r.FormValue("nps")); npsStr != "" {
		if score, err := strconv.Atoi(npsStr); err == nil && score >= 0 && score <= 10 {
			npsScore = score
		}
	}

	// Create survey response
	response := SurveyResponse{
		ID:                 id,
		Timestamp:          time.Now(),
		Likelihood:         strings.TrimSpace(r.FormValue("likelihood")),
		Tools:              r.Form["tools"], // Multiple checkbox values
		UseCase:            strings.TrimSpace(r.FormValue("use_case")),
		BusinessSector:     strings.TrimSpace(r.FormValue("business_sector")),
		EnterpriseInterest: strings.TrimSpace(r.FormValue("enterprise_interest")),
		Concerns:           strings.TrimSpace(r.FormValue("concerns")),
		FeatureRequest:     strings.TrimSpace(r.FormValue("feature_request")),
		NPS:                npsScore,
		Email:              strings.TrimSpace(r.FormValue("email")),
		Name:               strings.TrimSpace(r.FormValue("name")),
		Updates:            r.FormValue("updates") == "yes",
		IPAddress:          getClientIP(r),
	}

	// Save to file
	if err := saveSurveyToFile(response); err != nil {
		fmt.Printf("Failed to save survey response: %v\n", err)
		http.Error(w, "Failed to save feedback. Please try again later.", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Feedback Submitted - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script>
		<script>
			tailwind.config = {
				darkMode: 'class',
			}
		</script>
		<style>
			.theme-transition {
				transition: background-color 0.3s ease, color 0.3s ease, border-color 0.3s ease;
			}
		</style>
	</head>
	<body class="bg-gray-50 dark:bg-gray-900 theme-transition">
		<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
			<div class="bg-white dark:bg-gray-800 rounded-lg p-8 max-w-md mx-4 border border-gray-200 dark:border-gray-700 theme-transition">
				<div class="text-center">
					<div class="w-16 h-16 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg class="w-8 h-8 text-green-600 dark:text-green-400" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
						</svg>
					</div>
					<h3 class="text-xl font-semibold text-gray-800 dark:text-gray-100 mb-2">Thank You!</h3>
					<p class="text-gray-600 dark:text-gray-300 mb-6">Your feedback has been recorded and will help us improve ZeePass.</p>
					<button onclick="window.location.href='/'" class="bg-blue-600 hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600 text-white px-6 py-2 rounded-lg transition theme-transition">
						Back to Home
					</button>
				</div>
			</div>
		</div>

		<script>
			// Initialize theme from localStorage or system preference
			function initTheme() {
				const savedTheme = localStorage.getItem('theme');
				const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
				
				if (savedTheme === 'dark' || (!savedTheme && systemDark)) {
					document.documentElement.classList.add('dark');
				}
			}

			// Initialize on page load
			initTheme();

			// Auto-redirect after 5 seconds
			setTimeout(function() {
				window.location.href = '/';
			}, 5000);
		</script>
	</body>
	</html>
	`)
}

func saveSurveyToFile(response SurveyResponse) error {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Define file path
	filePath := "data/survey_responses.json"

	// Read existing responses
	var responses []SurveyResponse
	if fileData, err := os.ReadFile(filePath); err == nil {
		// File exists, unmarshal existing data
		if err := json.Unmarshal(fileData, &responses); err != nil {
			// If unmarshal fails, start with empty array
			responses = []SurveyResponse{}
		}
	}

	// Add new response to the array
	responses = append(responses, response)

	// Marshal the updated array
	jsonData, err := json.MarshalIndent(responses, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal survey responses: %v", err)
	}

	// Write the entire array back to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	// Also log to console for immediate visibility
	fmt.Printf("Survey response saved: ID=%s, Likelihood=%s, NPS=%d, Tools=%v, BusinessSector=%s, EnterpriseInterest=%s (Total responses: %d)\n", 
		response.ID, response.Likelihood, response.NPS, response.Tools, response.BusinessSector, response.EnterpriseInterest, len(responses))

	return nil
}

func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header (common in proxy/load balancer setups)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check for X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx] // Remove port
	}
	return ip
}