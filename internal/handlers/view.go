package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anazri/zeepass/internal/models"
	"github.com/anazri/zeepass/internal/services"
)

func ViewEncryptedHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := pathParts[2]

	log.Printf("[%s %s] Attempting to retrieve message with ID: %s", r.Method, r.RemoteAddr, id)
	data, err := services.GetEncryptedData(id)
	if err != nil {
		log.Printf("Failed to retrieve message ID %s: %v", id, err)
		html := `
		<!DOCTYPE html>
		<html><head><title>Message Not Found - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-red-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">Message Not Found</h2>
				<p class="text-gray-600 mb-6">This encrypted message does not exist or has expired.</p>
				<a href="/" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Go Home</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	if r.Method == http.MethodPost {
		// Re-fetch data for POST requests to ensure it's still available
		freshData, err := services.GetEncryptedData(id)
		if err != nil {
			log.Printf("Message %s no longer available for decryption: %v", id, err)
			html := `
			<!DOCTYPE html>
			<html><head><title>Message Not Found - ZeePass</title>
			<script src="https://cdn.tailwindcss.com"></script></head>
			<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
				<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
					<div class="text-red-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1-1z" clip-rule="evenodd"/></svg></div>
					<h2 class="text-2xl font-bold text-gray-800 mb-4">Message Not Found</h2>
					<p class="text-gray-600 mb-6">This encrypted message does not exist or has expired.</p>
					<a href="/" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Go Home</a>
				</div>
			</body></html>
			`
			w.Write([]byte(html))
			return
		}
		handleDecryptMessageWithData(w, r, id, freshData)
		return
	}

	if data.ExpiresAt != nil && time.Now().After(*data.ExpiresAt) {
		services.DeleteEncryptedData(id)
		html := `
		<!DOCTYPE html>
		<html><head><title>Message Expired - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-orange-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">Message Expired</h2>
				<p class="text-gray-600 mb-6">This encrypted message has expired and is no longer available.</p>
				<a href="/" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Go Home</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	if data.ViewCount >= data.MaxViews {
		services.DeleteEncryptedData(id)
		html := `
		<!DOCTYPE html>
		<html><head><title>Message No Longer Available - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-gray-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path d="M10 12a2 2 0 100-4 2 2 0 000 4z"/><path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">Message Already Viewed</h2>
				<p class="text-gray-600 mb-6">This message was configured to be viewed once and has already been accessed.</p>
				<a href="/" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Go Home</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	if data.PIN != "" {
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html><head>
			<title>Enter PIN - ZeePass</title>
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<script src="https://cdn.tailwindcss.com"></script>
		</head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-4 sm:p-8 rounded-lg shadow-md max-w-md w-full">
				<div class="text-center mb-4 sm:mb-6">
					<div class="w-12 h-12 sm:w-16 sm:h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
						<svg class="w-6 h-6 sm:w-8 sm:h-8 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd"/>
						</svg>
					</div>
					<h2 class="text-xl sm:text-2xl font-bold text-gray-800 mb-2">Protected Message</h2>
					<p class="text-sm sm:text-base text-gray-600">This message is protected with a PIN. Enter the PIN to view the content.</p>
				</div>
				<form method="POST">
					<div class="mb-4">
						<label class="block text-sm font-medium text-gray-700 mb-2">PIN</label>
						<input type="password" name="pin" required class="w-full px-3 py-3 sm:py-2 text-base sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none" placeholder="Enter PIN">
					</div>
					<button type="submit" class="w-full bg-blue-600 text-white py-3 sm:py-2 text-base sm:text-sm rounded-lg hover:bg-blue-700 transition font-medium">View Message</button>
				</form>
			</div>
		</body></html>
		`)
		w.Write([]byte(html))
		return
	}

	// Show a confirmation form for non-PIN protected messages
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html><head>
		<title>View Message - ZeePass</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
		<div class="bg-white p-4 sm:p-8 rounded-lg shadow-md max-w-md w-full">
			<div class="text-center mb-4 sm:mb-6">
				<div class="w-12 h-12 sm:w-16 sm:h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
					<svg class="w-6 h-6 sm:w-8 sm:h-8 text-green-600" fill="currentColor" viewBox="0 0 20 20">
						<path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
					</svg>
				</div>
				<h2 class="text-xl sm:text-2xl font-bold text-gray-800 mb-2">Encrypted Message</h2>
				<p class="text-sm sm:text-base text-gray-600">Click below to view this encrypted message.</p>
				%s
			</div>
			<form method="POST">
				<button type="submit" class="w-full bg-green-600 text-white py-3 sm:py-2 text-base sm:text-sm rounded-lg hover:bg-green-700 transition font-medium">View Message</button>
			</form>
		</div>
	</body></html>
	`, getWarningMessage(data))
	w.Write([]byte(html))
}

func handleDecryptMessageWithData(w http.ResponseWriter, r *http.Request, id string, data *models.EncryptedData) {
	pin := r.FormValue("pin")

	if data.PIN != "" && services.HashPIN(pin) != data.PIN {
		html := `
		<!DOCTYPE html>
		<html><head><title>Invalid PIN - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-red-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">Invalid PIN</h2>
				<p class="text-gray-600 mb-6">The PIN you entered is incorrect.</p>
				<a href="javascript:history.back()" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Try Again</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	showDecryptedMessageWithData(w, r, id, data)
}

func showDecryptedMessageWithData(w http.ResponseWriter, r *http.Request, id string, data *models.EncryptedData) {
	data.ViewCount++

	if data.ViewCount >= data.MaxViews {
		err := services.DeleteEncryptedData(id)
		if err != nil {
			log.Printf("Error deleting message after max views: %v", err)
		}
	} else {
		err := services.StoreEncryptedData(id, data)
		if err != nil {
			log.Printf("Error updating view count in Redis: %v", err)
		}
	}

	decryptedText, err := services.Decrypt(data.Content, services.GetEncryptionKey())
	if err != nil {
		http.Error(w, "Error decrypting message", http.StatusInternalServerError)
		return
	}

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html><head>
		<title>Encrypted Message - ZeePass</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body class="bg-gray-50 min-h-screen py-4 sm:py-8">
		<div class="max-w-4xl mx-auto px-2 sm:px-4">
			<div class="bg-white rounded-lg shadow-md overflow-hidden">
				<div class="bg-green-500 text-white p-3 sm:p-4">
					<div class="flex items-center space-x-2">
						<svg class="w-5 h-5 sm:w-6 sm:h-6" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 616 0z" clip-rule="evenodd"/></svg>
						<h1 class="text-lg sm:text-xl font-bold">Decrypted Message</h1>
					</div>
				</div>
				<div class="p-3 sm:p-6">
					<div class="mb-4">
						<label class="block text-sm font-medium text-gray-700 mb-2">Message Content</label>
						<div class="bg-gray-50 p-2 sm:p-4 rounded-lg border">
							<pre class="whitespace-pre-wrap text-sm sm:text-base text-gray-800 overflow-x-auto break-words">%s</pre>
						</div>
					</div>
					%s
					<div class="flex flex-col sm:flex-row sm:justify-between items-center mt-4 sm:mt-6 space-y-2 sm:space-y-0 sm:space-x-4">
						<button onclick="copyMessage()" class="w-full sm:w-auto bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition text-center">Copy Message</button>
						<a href="/" class="w-full sm:w-auto bg-gray-600 text-white px-4 py-2 rounded-lg hover:bg-gray-700 transition text-center">Create New Message</a>
					</div>
				</div>
			</div>
		</div>
		<script>
			function copyMessage() {
				const messageText = %s;
				navigator.clipboard.writeText(messageText).then(() => {
					alert('Message copied to clipboard!');
				}).catch(() => {
					alert('Failed to copy message. Please select and copy manually.');
				});
			}
		</script>
	</body></html>
	`, decryptedText, getWarningMessage(data), strconv.Quote(decryptedText))

	w.Write([]byte(html))
}

func getWarningMessage(data *models.EncryptedData) string {
	if data.MaxViews == 1 {
		return `<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">⚠️ <strong>Warning:</strong> This message will be permanently deleted after viewing.</div>`
	}
	return ""
}

func getFileWarningMessage(data *models.EncryptedFileData) string {
	if data.MaxViews == 1 {
		return `<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">⚠️ <strong>Warning:</strong> This file will be permanently deleted after downloading.</div>`
	}
	return ""
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

func ViewEncryptedFileHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.NotFound(w, r)
		return
	}
	id := pathParts[2]

	log.Printf("[%s %s] Attempting to retrieve file with ID: %s", r.Method, r.RemoteAddr, id)
	data, err := services.GetEncryptedFileData(id)
	if err != nil {
		log.Printf("Failed to retrieve file ID %s: %v", id, err)
		html := `
		<!DOCTYPE html>
		<html><head><title>File Not Found - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-red-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">File Not Found</h2>
				<p class="text-gray-600 mb-6">This encrypted file does not exist or has expired.</p>
				<a href="/" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Go Home</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	if r.Method == http.MethodPost {
		handleDecryptFileWithData(w, r, id, data)
		return
	}

	if data.ExpiresAt != nil && time.Now().After(*data.ExpiresAt) {
		services.DeleteEncryptedFileData(id)
		html := `
		<!DOCTYPE html>
		<html><head><title>File Expired - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-orange-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">File Expired</h2>
				<p class="text-gray-600 mb-6">This encrypted file has expired and is no longer available.</p>
				<a href="/" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Go Home</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	if data.ViewCount >= data.MaxViews {
		services.DeleteEncryptedFileData(id)
		html := `
		<!DOCTYPE html>
		<html><head><title>File No Longer Available - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-gray-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path d="M10 12a2 2 0 100-4 2 2 0 000 4z"/><path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">File Already Downloaded</h2>
				<p class="text-gray-600 mb-6">This file was configured to be downloaded once and has already been accessed.</p>
				<a href="/" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Go Home</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	if data.PIN != "" {
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html><head>
			<title>Enter PIN - ZeePass</title>
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<script src="https://cdn.tailwindcss.com"></script>
		</head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-4 sm:p-8 rounded-lg shadow-md max-w-md w-full">
				<div class="text-center mb-4 sm:mb-6">
					<div class="w-12 h-12 sm:w-16 sm:h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
						<svg class="w-6 h-6 sm:w-8 sm:h-8 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd"/>
						</svg>
					</div>
					<h2 class="text-xl sm:text-2xl font-bold text-gray-800 mb-2">Protected File</h2>
					<p class="text-sm sm:text-base text-gray-600">This file is protected with a PIN. Enter the PIN to download the file.</p>
				</div>
				<form method="POST">
					<div class="mb-4">
						<label class="block text-sm font-medium text-gray-700 mb-2">PIN</label>
						<input type="password" name="pin" required class="w-full px-3 py-3 sm:py-2 text-base sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none" placeholder="Enter PIN">
					</div>
					<button type="submit" class="w-full bg-blue-600 text-white py-3 sm:py-2 text-base sm:text-sm rounded-lg hover:bg-blue-700 transition font-medium">Download File</button>
				</form>
			</div>
		</body></html>
		`)
		w.Write([]byte(html))
		return
	}

	// Show a confirmation form for non-PIN protected files
	fileSize := formatFileSize(data.FileSize)
	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html><head>
		<title>Download File - ZeePass</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<script src="https://cdn.tailwindcss.com"></script>
	</head>
	<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
		<div class="bg-white p-4 sm:p-8 rounded-lg shadow-md max-w-md w-full">
			<div class="text-center mb-4 sm:mb-6">
				<div class="w-12 h-12 sm:w-16 sm:h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
					<svg class="w-6 h-6 sm:w-8 sm:h-8 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd"/>
					</svg>
				</div>
				<h2 class="text-xl sm:text-2xl font-bold text-gray-800 mb-2">Encrypted File</h2>
				<p class="text-sm sm:text-base text-gray-600">Click below to download this encrypted file.</p>
				<div class="bg-gray-50 p-3 rounded border mt-4 mb-4">
					<p><strong>File Name:</strong> %s</p>
					<p><strong>File Size:</strong> %s</p>
				</div>
				%s
			</div>
			<form method="POST">
				<button type="submit" class="w-full bg-blue-600 text-white py-3 sm:py-2 text-base sm:text-sm rounded-lg hover:bg-blue-700 transition font-medium">Download File</button>
			</form>
		</div>
	</body></html>
	`, data.FileName, fileSize, getFileWarningMessage(data))
	w.Write([]byte(html))
}

func handleDecryptFileWithData(w http.ResponseWriter, r *http.Request, id string, data *models.EncryptedFileData) {
	pin := r.FormValue("pin")

	if data.PIN != "" && services.HashPIN(pin) != data.PIN {
		html := `
		<!DOCTYPE html>
		<html><head><title>Invalid PIN - ZeePass</title>
		<script src="https://cdn.tailwindcss.com"></script></head>
		<body class="bg-gray-50 flex items-center justify-center min-h-screen p-4">
			<div class="bg-white p-8 rounded-lg shadow-md text-center max-w-md">
				<div class="text-red-500 mb-4"><svg class="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/></svg></div>
				<h2 class="text-2xl font-bold text-gray-800 mb-4">Invalid PIN</h2>
				<p class="text-gray-600 mb-6">The PIN you entered is incorrect.</p>
				<a href="javascript:history.back()" class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition">Try Again</a>
			</div>
		</body></html>
		`
		w.Write([]byte(html))
		return
	}

	downloadDecryptedFileWithData(w, r, id, data)
}

func downloadDecryptedFileWithData(w http.ResponseWriter, r *http.Request, id string, data *models.EncryptedFileData) {
	data.ViewCount++

	if data.ViewCount >= data.MaxViews {
		err := services.DeleteEncryptedFileData(id)
		if err != nil {
			log.Printf("Error deleting file after max views: %v", err)
		}
	} else {
		err := services.StoreEncryptedFileData(id, data)
		if err != nil {
			log.Printf("Error updating view count in Redis: %v", err)
		}
	}

	decryptedData, err := services.DecryptFile(data.Content, services.GetEncryptionKey())
	if err != nil {
		http.Error(w, "Error decrypting file", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Type", data.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", data.FileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(decryptedData)))

	// Write file data
	w.Write(decryptedData)
}