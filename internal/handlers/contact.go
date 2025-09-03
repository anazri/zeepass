package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ContactForm struct {
	Name        string `form:"name"`
	Email       string `form:"email"`
	Company     string `form:"company"`
	InquiryType string `form:"inquiry_type"`
	Message     string `form:"message"`
	Website     string `form:"website"` // Honeypot field
}

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// Rate limiting
var (
	rateLimiter = make(map[string][]time.Time)
	rateMutex   = sync.RWMutex{}
	maxRequests = 3                // Max requests per time window
	timeWindow  = 10 * time.Minute // Time window for rate limiting
)

// Spam keywords to filter
var spamKeywords = []string{
	"viagra", "casino", "lottery", "bitcoin", "crypto", "investment",
	"make money", "get rich", "buy now", "click here", "free money",
	"guaranteed", "no risk", "limited time", "act now", "cheap",
	"discount", "save money", "earn cash", "work from home",
}

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
func HandleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Rate limiting check
	clientIP := getClientIP(r)
	if !checkRateLimit(clientIP) {
		http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	form := ContactForm{
		Name:        strings.TrimSpace(r.FormValue("name")),
		Email:       strings.TrimSpace(r.FormValue("email")),
		Company:     strings.TrimSpace(r.FormValue("company")),
		InquiryType: strings.TrimSpace(r.FormValue("inquiry_type")),
		Message:     strings.TrimSpace(r.FormValue("message")),
		Website:     strings.TrimSpace(r.FormValue("website")), // Honeypot
	}

	// Honeypot check - if filled, it's likely spam
	if form.Website != "" {
		http.Error(w, "Spam detected", http.StatusBadRequest)
		return
	}

	// Enhanced validation
	if form.Name == "" || form.Email == "" || form.Message == "" {
		http.Error(w, "Name, email, and message are required", http.StatusBadRequest)
		return
	}

	// Length validation
	if len(form.Name) > 100 || len(form.Email) > 255 || len(form.Company) > 100 {
		http.Error(w, "Input too long", http.StatusBadRequest)
		return
	}

	if len(form.Message) < 10 || len(form.Message) > 5000 {
		http.Error(w, "Message must be between 10 and 5000 characters", http.StatusBadRequest)
		return
	}

	// Email validation
	if !isValidEmail(form.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Spam content filtering
	if containsSpam(form.Message) || containsSpam(form.Name) {
		http.Error(w, "Message contains inappropriate content", http.StatusBadRequest)
		return
	}

	// Verify CAPTCHA (if enabled)
	recaptchaResponse := strings.TrimSpace(r.FormValue("g-recaptcha-response"))
	if !verifyCaptcha(recaptchaResponse, r.RemoteAddr) {
		http.Error(w, "CAPTCHA verification failed. Please try again.", http.StatusBadRequest)
		return
	}

	// Send email
	if err := sendContactEmail(form); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		http.Error(w, "Failed to send message. Please try again later or contact us directly.", http.StatusInternalServerError)
		return
	}

	// Return success response that matches the app's design system
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en" class="scroll-smooth">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Message Sent - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script>
		<script>
			tailwind.config = {
				darkMode: 'class',
				theme: {
					extend: {
						animation: {
							'fade-in': 'fadeIn 0.5s ease-in-out',
							'slide-up': 'slideUp 0.3s ease-out',
						}
					}
				}
			}
		</script>
		<style>
			.theme-transition {
				transition: background-color 0.3s ease, color 0.3s ease, border-color 0.3s ease;
			}
			@keyframes fadeIn {
				from { opacity: 0; }
				to { opacity: 1; }
			}
			@keyframes slideUp {
				from { transform: translateY(20px); opacity: 0; }
				to { transform: translateY(0); opacity: 1; }
			}
		</style>
		<script>
			// Apply saved theme on page load
			if (localStorage.getItem('darkMode') === 'true' || (!localStorage.getItem('darkMode') && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
				document.documentElement.classList.add('dark');
			}
		</script>
	</head>
	<body class="bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100 theme-transition">
		<div class="min-h-screen flex items-center justify-center px-4 py-8">
			<div class="max-w-md w-full">
				<div class="bg-white dark:bg-gray-800 rounded-2xl p-8 shadow-lg border border-gray-200 dark:border-gray-700 theme-transition animate-slide-up">
					<div class="text-center">
						<div class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
							<svg class="w-10 h-10 text-green-600 dark:text-green-400" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
							</svg>
						</div>
						<h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100 mb-3">Message Sent Successfully!</h1>
						<p class="text-gray-600 dark:text-gray-400 mb-8 leading-relaxed">
							Thank you for reaching out to us. We've received your message and will get back to you within 24 hours.
						</p>
						<div class="space-y-3">
							<button onclick="window.location.href='/'" 
									class="w-full bg-blue-600 hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-600 text-white py-3 px-6 rounded-lg font-semibold transition theme-transition">
								Return to Homepage
							</button>
							<button onclick="window.history.back()" 
									class="w-full bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 py-3 px-6 rounded-lg font-medium transition theme-transition">
								Go Back
							</button>
						</div>
					</div>
				</div>
				<div class="text-center mt-6">
					<p class="text-sm text-gray-500 dark:text-gray-400">
						Need immediate assistance? Contact us at 
						<a href="mailto:contact@moonkite.io" class="text-blue-600 dark:text-blue-400 hover:underline">contact@moonkite.io</a>
					</p>
				</div>
			</div>
		</div>
	</body>
	</html>
	`)
}

type ResendEmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

func sendContactEmail(form ContactForm) error {
	// Get Resend configuration from environment variables
	resendAPIKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	toEmail := os.Getenv("CONTACT_EMAIL")

	// Use defaults if not set
	if fromEmail == "" {
		fromEmail = "noreply@resend.dev"
	}
	if toEmail == "" {
		toEmail = "contact@moonkite.io"
	}

	// If Resend API key is not configured, log the message instead
	if resendAPIKey == "" {
		fmt.Printf("Resend API key not configured. Contact form submission:\n")
		fmt.Printf("Name: %s\n", form.Name)
		fmt.Printf("Email: %s\n", form.Email)
		fmt.Printf("Company: %s\n", form.Company)
		fmt.Printf("Inquiry Type: %s\n", form.InquiryType)
		fmt.Printf("Message: %s\n", form.Message)
		fmt.Printf("Timestamp: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil
	}

	// Create email content
	subject := fmt.Sprintf("ZeePass Contact Form: %s", getInquiryTypeLabel(form.InquiryType))

	htmlBody := fmt.Sprintf(`
		<h2>New contact form submission from ZeePass website</h2>
		
		<p><strong>Name:</strong> %s</p>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Company:</strong> %s</p>
		<p><strong>Inquiry Type:</strong> %s</p>
		<p><strong>Timestamp:</strong> %s</p>
		
		<h3>Message:</h3>
		<p>%s</p>
		
		<hr>
		<p><em>This email was sent automatically from the ZeePass contact form.</em></p>
		`,
		form.Name,
		form.Email,
		form.Company,
		getInquiryTypeLabel(form.InquiryType),
		time.Now().Format("2006-01-02 15:04:05"),
		strings.ReplaceAll(form.Message, "\n", "<br>"),
	)

	// Create Resend email request
	emailReq := ResendEmailRequest{
		From:    fromEmail,
		To:      toEmail,
		Subject: subject,
		HTML:    htmlBody,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	// Debug log the JSON payload
	fmt.Printf("Sending email with payload: %s\n", string(jsonData))

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+resendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resend API returned status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Email sent successfully via Resend API\n")
	return nil
}

func getInquiryTypeLabel(inquiryType string) string {
	switch inquiryType {
	case "cloud":
		return "Cloud Hosting Setup"
	case "onpremise":
		return "On-Premise Installation"
	case "support":
		return "Technical Support"
	case "consultation":
		return "Consultation & Planning"
	case "custom":
		return "Custom Development"
	case "question":
		return "General Question"
	case "other":
		return "Other"
	default:
		return "General Inquiry"
	}
}

func isValidEmail(email string) bool {
	// Enhanced email validation with regex
	if len(email) > 255 {
		return false
	}
	return emailRegex.MatchString(email)
}

// Rate limiting functions

func checkRateLimit(clientIP string) bool {
	rateMutex.Lock()
	defer rateMutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-timeWindow)
	
	// Get existing requests for this IP
	requests := rateLimiter[clientIP]
	
	// Filter out old requests
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check if under limit
	if len(validRequests) >= maxRequests {
		return false
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	rateLimiter[clientIP] = validRequests
	
	return true
}

// Spam filtering function
func containsSpam(text string) bool {
	lowerText := strings.ToLower(text)
	
	// Check for spam keywords
	for _, keyword := range spamKeywords {
		if strings.Contains(lowerText, keyword) {
			return true
		}
	}
	
	// Check for excessive links (more than 2 URLs)
	urlCount := strings.Count(lowerText, "http://") + strings.Count(lowerText, "https://") + strings.Count(lowerText, "www.")
	if urlCount > 2 {
		return true
	}
	
	// Check for excessive repetitive characters
	repeatPattern := regexp.MustCompile(`(.)\1{4,}`) // 5 or more repeated characters
	if repeatPattern.MatchString(text) {
		return true
	}
	
	return false
}


func verifyCaptcha(response, remoteIP string) bool {
	recaptchaSecret := os.Getenv("RECAPTCHA_SECRET_KEY")
	if recaptchaSecret == "" {
		// If no secret key is configured, skip CAPTCHA verification in development
		fmt.Println("Warning: RECAPTCHA_SECRET_KEY not configured, skipping CAPTCHA verification")
		return true
	}

	// If no response provided and CAPTCHA is enabled, fail verification
	if response == "" {
		fmt.Println("CAPTCHA response is required but not provided")
		return false
	}

	// Prepare form data for Google reCAPTCHA API
	formData := url.Values{}
	formData.Set("secret", recaptchaSecret)
	formData.Set("response", response)
	formData.Set("remoteip", remoteIP)

	// Make request to Google reCAPTCHA API
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", formData)
	if err != nil {
		fmt.Printf("Error verifying CAPTCHA: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading CAPTCHA response: %v\n", err)
		return false
	}

	// Parse JSON response
	var recaptchaResp RecaptchaResponse
	if err := json.Unmarshal(body, &recaptchaResp); err != nil {
		fmt.Printf("Error parsing CAPTCHA response: %v\n", err)
		return false
	}

	// Check if verification was successful
	if !recaptchaResp.Success {
		fmt.Printf("CAPTCHA verification failed: %v\n", recaptchaResp.ErrorCodes)
		return false
	}

	// For reCAPTCHA v3, you can also check the score (0.0-1.0, higher is more likely to be human)
	// Uncomment and adjust threshold as needed:
	// if recaptchaResp.Score < 0.5 {
	//     fmt.Printf("CAPTCHA score too low: %f\n", recaptchaResp.Score)
	//     return false
	// }

	fmt.Printf("CAPTCHA verification successful (score: %f)\n", recaptchaResp.Score)
	return true
}
