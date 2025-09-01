package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ContactForm struct {
	Name        string `form:"name"`
	Email       string `form:"email"`
	Company     string `form:"company"`
	InquiryType string `form:"inquiry_type"`
	Message     string `form:"message"`
}

func HandleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	}

	// Basic validation
	if form.Name == "" || form.Email == "" || form.Message == "" {
		http.Error(w, "Name, email, and message are required", http.StatusBadRequest)
		return
	}

	if !isValidEmail(form.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Send email
	if err := sendContactEmail(form); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		http.Error(w, "Failed to send message. Please try again later or contact us directly.", http.StatusInternalServerError)
		return
	}

	// Return success response (you could redirect to a thank you page)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-8 max-w-md mx-4">
			<div class="text-center">
				<div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
					<svg class="w-8 h-8 text-green-600" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
					</svg>
				</div>
				<h3 class="text-xl font-semibold text-gray-800 mb-2">Message Sent!</h3>
				<p class="text-gray-600 mb-6">Thank you for contacting us. We'll get back to you within 24 hours.</p>
				<button onclick="window.location.reload()" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">
					Close
				</button>
			</div>
		</div>
	</div>
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
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}
