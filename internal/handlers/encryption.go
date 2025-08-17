package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anazri/zeepass/internal/models"
	"github.com/anazri/zeepass/internal/services"
)

func EncryptTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Error parsing form data</div>`)
		w.Write([]byte(responseHTML))
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	pin := r.FormValue("pin")
	lifetime := r.FormValue("lifetime")

	if text == "" {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Please enter some text to encrypt</div>`)
		w.Write([]byte(responseHTML))
		return
	}

	id := services.GenerateID()

	encryptedText, err := services.Encrypt(text, services.GetEncryptionKey())
	if err != nil {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Error encrypting text: %v</div>`, err)
		w.Write([]byte(responseHTML))
		return
	}

	hashedPIN := ""
	if pin != "" {
		hashedPIN = services.HashPIN(pin)
	}

	var expiresAt *time.Time
	maxViews := 1
	switch lifetime {
	case "1h":
		expiry := time.Now().Add(time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "24h":
		expiry := time.Now().Add(24 * time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "7d":
		expiry := time.Now().Add(7 * 24 * time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "30d":
		expiry := time.Now().Add(30 * 24 * time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "never":
		maxViews = 999999
	default:
		maxViews = 1
	}

	encData := &models.EncryptedData{
		ID:        id,
		Content:   encryptedText,
		PIN:       hashedPIN,
		Lifetime:  lifetime,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		ViewCount: 0,
		MaxViews:  maxViews,
	}

	err = services.StoreEncryptedData(id, encData)
	if err != nil {
		log.Printf("Error storing encrypted data for ID %s: %v", id, err)
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Error storing encrypted data: %v</div>`, err)
		w.Write([]byte(responseHTML))
		return
	}
	log.Printf("Successfully stored encrypted data for ID: %s", id)

	viewURL := fmt.Sprintf("http://localhost:8080/view/%s", id)

	responseHTML := fmt.Sprintf(`
		<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
			✅ Text encrypted successfully!
		</div>
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<h3 class="text-lg font-semibold text-gray-800 mb-4">Share this encrypted message</h3>
			<div class="mb-4">
				<label class="block text-sm font-medium text-gray-700 mb-2">Secure Link</label>
				<div class="flex">
					<input type="text" value="%s" readonly class="flex-1 px-3 py-2 border border-gray-300 rounded-l-lg bg-gray-50 text-sm" id="shareURL">
					<button onclick="copyToClipboard()" class="px-4 py-2 bg-blue-600 text-white rounded-r-lg hover:bg-blue-700 transition">
						<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
							<path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z"/>
							<path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2 3 3 0 01-3 3H9a3 3 0 01-3-3z"/>
						</svg>
					</button>
				</div>
			</div>
			<div class="text-sm text-gray-600">
				<p><strong>Lifetime:</strong> %s</p>
				%s
				<p class="mt-2 text-amber-600">⚠️ This link will expire according to the lifetime settings. Save it securely.</p>
			</div>
		</div>
		<script>
			function copyToClipboard() {
				const urlInput = document.getElementById('shareURL');
				urlInput.select();
				document.execCommand('copy');
				alert('Link copied to clipboard!');
			}
		</script>
	`, viewURL, getLifetimeDisplay(lifetime), getPINDisplay(pin))

	w.Write([]byte(responseHTML))
}

func getLifetimeDisplay(lifetime string) string {
	switch lifetime {
	case "1h":
		return "1 Hour"
	case "24h":
		return "24 Hours"
	case "7d":
		return "7 Days"
	case "30d":
		return "30 Days"
	case "never":
		return "Never expires"
	default:
		return "Once received"
	}
}

func getPINDisplay(pin string) string {
	if pin != "" {
		return "<p><strong>PIN Protection:</strong> Enabled</p>"
	}
	return "<p><strong>PIN Protection:</strong> Not set</p>"
}

func EncryptFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 10MB limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Error parsing form data</div>`)
		w.Write([]byte(responseHTML))
		return
	}

	// Get form values
	pin := r.FormValue("pin")
	lifetime := r.FormValue("lifetime")

	// Get uploaded file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Please select a file to encrypt</div>`)
		w.Write([]byte(responseHTML))
		return
	}
	defer file.Close()

	// Check file size (10MB limit)
	if fileHeader.Size > 10<<20 {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">File size must be less than 10MB</div>`)
		w.Write([]byte(responseHTML))
		return
	}

	// Read file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Error reading file: %v</div>`, err)
		w.Write([]byte(responseHTML))
		return
	}

	// Generate ID for the encrypted file
	id := services.GenerateID()

	// Encrypt file data
	encryptedData, err := services.EncryptFile(fileData, services.GetEncryptionKey())
	if err != nil {
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Error encrypting file: %v</div>`, err)
		w.Write([]byte(responseHTML))
		return
	}

	// Hash PIN if provided
	hashedPIN := ""
	if pin != "" {
		hashedPIN = services.HashPIN(pin)
	}

	// Set expiration time and max views based on lifetime
	var expiresAt *time.Time
	maxViews := 1
	switch lifetime {
	case "1h":
		expiry := time.Now().Add(time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "24h":
		expiry := time.Now().Add(24 * time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "7d":
		expiry := time.Now().Add(7 * 24 * time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "30d":
		expiry := time.Now().Add(30 * 24 * time.Hour)
		expiresAt = &expiry
		maxViews = 999999
	case "never":
		maxViews = 999999
	default:
		maxViews = 1
	}

	// Create encrypted file data struct
	encFileData := &models.EncryptedFileData{
		ID:        id,
		Content:   encryptedData,
		FileName:  fileHeader.Filename,
		FileSize:  fileHeader.Size,
		MimeType:  fileHeader.Header.Get("Content-Type"),
		PIN:       hashedPIN,
		Lifetime:  lifetime,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		ViewCount: 0,
		MaxViews:  maxViews,
	}

	// Store encrypted file data
	err = services.StoreEncryptedFileData(id, encFileData)
	if err != nil {
		log.Printf("Error storing encrypted file data for ID %s: %v", id, err)
		responseHTML := fmt.Sprintf(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">Error storing encrypted file: %v</div>`, err)
		w.Write([]byte(responseHTML))
		return
	}
	log.Printf("Successfully stored encrypted file data for ID: %s", id)

	// Generate view URL
	viewURL := fmt.Sprintf("http://localhost:8080/view-file/%s", id)

	// Calculate file size in human-readable format
	fileSize := formatFileSize(fileHeader.Size)

	// Generate success response HTML
	responseHTML := fmt.Sprintf(`
		<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
			✅ File encrypted successfully!
		</div>
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<h3 class="text-lg font-semibold text-gray-800 mb-4">Share this encrypted file</h3>
			<div class="mb-4">
				<label class="block text-sm font-medium text-gray-700 mb-2">File Details</label>
				<div class="bg-gray-50 p-3 rounded border">
					<p><strong>File Name:</strong> %s</p>
					<p><strong>File Size:</strong> %s</p>
				</div>
			</div>
			<div class="mb-4">
				<label class="block text-sm font-medium text-gray-700 mb-2">Secure Link</label>
				<div class="flex">
					<input type="text" value="%s" readonly class="flex-1 px-3 py-2 border border-gray-300 rounded-l-lg bg-gray-50 text-sm" id="shareURL">
					<button onclick="copyToClipboard()" class="px-4 py-2 bg-blue-600 text-white rounded-r-lg hover:bg-blue-700 transition">
						<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
							<path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z"/>
							<path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2 3 3 0 01-3 3H9a3 3 0 01-3-3z"/>
						</svg>
					</button>
				</div>
			</div>
			<div class="text-sm text-gray-600">
				<p><strong>Lifetime:</strong> %s</p>
				%s
				<p class="mt-2 text-amber-600">⚠️ This link will expire according to the lifetime settings. Save it securely.</p>
			</div>
		</div>
		<script>
			function copyToClipboard() {
				const urlInput = document.getElementById('shareURL');
				urlInput.select();
				document.execCommand('copy');
				alert('Link copied to clipboard!');
			}
		</script>
	`, fileHeader.Filename, fileSize, viewURL, getLifetimeDisplay(lifetime), getPINDisplay(pin))

	w.Write([]byte(responseHTML))
}

func formatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 Bytes"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d Bytes", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}