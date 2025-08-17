# ZeePass - Secure Encryption & Crypto Tools

🔐 **ZeePass** is a comprehensive web application providing end-to-end encryption tools for text, files, chat, password generation, SSH keys, and Base64 encoding/decoding. Built with Go and HTMX for a modern, secure, and user-friendly experience.

![ZeePass](https://img.shields.io/badge/Go-1.24.2-blue) ![License](https://img.shields.io/badge/License-MIT-green) ![Security](https://img.shields.io/badge/Security-AES--256--GCM-red)

## 🌟 Features

### 🔒 **Text Encryption**
- **AES-256-GCM encryption** for maximum security
- **PIN protection** with SHA-256 hashing
- **Configurable lifetime**: Once-read, 1 hour, 24 hours, 7 days, 30 days, or never expires
- **Auto-destruction** after reading (for once-read messages)
- **Secure sharing** via unique URLs

### 📄 **File Encryption**
- **Encrypt any file type** up to 10MB
- **AES-256-GCM encryption** with same security features as text
- **File metadata protection** (filename, size, MIME type)
- **Secure download** with automatic cleanup
- Support for **PIN protection and lifetime management**

### 💬 **Chat Encryption**
- **Real-time encrypted chat** via WebSockets
- **End-to-end encryption** - messages encrypted before transmission
- **Auto-expiring messages** with configurable lifetime
- **Redis-backed storage** for scalability
- **No message logging** - everything is encrypted

### 🔑 **Password Generator**
- **Multiple password types**:
  - Random passwords with customizable character sets
  - Memorable passwords using word combinations
  - PIN codes for secure access
- **Strength analysis** (weak/medium/strong)
- **Configurable length** (4-64 characters)
- **Character set options**: uppercase, lowercase, numbers, symbols

### 🔐 **SSH Key Generator**
- **Multiple key types**: RSA, Ed25519, ECDSA
- **Key length options**: 
  - RSA: 2048, 3072, 4096 bits
  - ECDSA: 256, 384, 521 bits
  - Ed25519: 256 bits (fixed)
- **Passphrase protection** with AES-256 encryption
- **Custom comments** for key identification
- **Industry-standard formats** (PEM, OpenSSH)

### 📋 **Base64 Tools**
- **Encode/Decode text** to/from Base64
- **Real-time processing** with HTMX
- **Error handling** for invalid Base64 input
- **Clean, intuitive interface**

### 🎨 **User Interface & Experience**
- **Modern Landing Page** with comprehensive feature overview
- **Dark/Light Theme Toggle** with system preference detection
- **Responsive Design** optimized for all device sizes
- **About Us Section** featuring:
  - Our Mission: Making advanced cryptography accessible to everyone
  - Our Expertise: Built by security experts and cryptography specialists  
  - Our Commitment: Continuous innovation and transparency in security
- **Professional Contact Forms** for enterprise deployment inquiries
- **Smooth Animations** and glassmorphism design effects

## 🏗️ Architecture

### **Backend (Go)**
```
cmd/server/          # Application entry point
internal/
├── handlers/        # HTTP request handlers
├── models/          # Data structures
└── services/        # Business logic
    ├── crypto.go    # Encryption/decryption
    ├── storage.go   # Redis data persistence
    ├── password.go  # Password generation
    ├── sshkey.go    # SSH key generation
    └── chat.go      # Real-time chat
```

### **Frontend (HTMX + TailwindCSS + JavaScript)**
```
templates/           # HTML templates with responsive design
├── index.html      # Landing page with dark mode support
├── base64.html     # Base64 encoding/decoding tool
├── chat-encryption.html  # Real-time encrypted chat
├── file-encryption.html  # File upload and encryption
├── password-generator.html  # Password generation tool
├── ssh-key.html    # SSH key generation tool
└── text-encryption.html   # Text encryption and sharing
```

**Frontend Features:**
- **Dark/Light Theme Toggle** with localStorage persistence
- **Responsive Design** optimized for mobile and desktop
- **System Theme Detection** (follows OS preferences)
- **Smooth Animations** and transitions
- **Modern UI Components** with glassmorphism effects

### **Storage (Redis)**
- **Encrypted data storage** with automatic TTL
- **Chat message persistence**
- **View count tracking**
- **Automatic cleanup** of expired content

## 🚀 Quick Start

### Prerequisites
- **Go 1.24.2+**
- **Redis Server** (for data persistence)
- **Git** (for cloning)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/anazri/zeepass.git
   cd zeepass
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start Redis server**
   ```bash
   # On macOS with Homebrew
   brew services start redis
   
   # On Ubuntu/Debian
   sudo systemctl start redis-server
   
   # On Windows (with Redis for Windows)
   redis-server
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the application**
   ```
   Open your browser and navigate to: http://localhost:8080
   ```

## 🔧 Configuration

### **Redis Configuration**
Edit `internal/services/storage.go` to configure Redis connection:
```go
rdb = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "your-redis-password",
    DB:       0,
})
```

### **Encryption Key**
**⚠️ IMPORTANT**: Change the default encryption key in `internal/services/crypto.go`:
```go
var encryptionKey = []byte("your-32-byte-encryption-key-here")
```
Use a cryptographically secure 32-byte key in production.

## 🛡️ Security Features

- **AES-256-GCM Encryption**: Military-grade encryption for all data
- **Zero-Knowledge Architecture**: Server never sees unencrypted data
- **SHA-256 PIN Hashing**: Secure PIN protection
- **Automatic Expiration**: Time-based and view-count-based expiry
- **Secure Random Generation**: Cryptographically secure random number generation
- **TLS-Ready**: Designed for HTTPS deployment

## 🌐 Deployment

### **Development**
```bash
# Run directly
go run cmd/server/main.go

# Or use Air for live reloading (recommended)
air
```

### **Production Build**
```bash
go build -o zeepass cmd/server/main.go
./zeepass
```

### **Docker Deployment**

```dockerfile
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o zeepass cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/zeepass .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./zeepass"]
```

### **Environment Variables**
- `REDIS_URL`: Redis connection string
- `ENCRYPTION_KEY`: 32-byte encryption key (base64 encoded)
- `PORT`: Server port (default: 8080)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Nazri Abdullah**  
Senior Software Engineer & Cybersecurity Expert  
[LinkedIn](https://www.linkedin.com/in/nazri-bin-abdullah-a75b40268/) | [GitHub](https://github.com/anazri)

---

## 🛠️ Tech Stack

- **Backend**: Go 1.24.2
- **Frontend**: HTMX, TailwindCSS, Vanilla JavaScript
- **Database**: Redis
- **Encryption**: AES-256-GCM
- **WebSockets**: Gorilla WebSocket
- **Development**: Air (live reloading)
- **Deployment**: Docker-ready

## 📊 Project Status

- ✅ **Text Encryption** - Complete
- ✅ **File Encryption** - Complete  
- ✅ **Chat Encryption** - Complete
- ✅ **Password Generator** - Complete
- ✅ **SSH Key Generator** - Complete
- ✅ **Base64 Tools** - Complete
- ✅ **Dark/Light Theme** - Complete
- ✅ **Responsive Design** - Complete
- 🔄 **User Authentication** - Planned
- 🔄 **API Endpoints** - Planned
- 🔄 **Mobile App** - Planned

---

⭐ **Star this repository if you find it useful!**

🐛 **Found a bug?** [Report it here](https://github.com/anazri/zeepass/issues)

💡 **Have a suggestion?** [Let us know!](https://github.com/anazri/zeepass/discussions)