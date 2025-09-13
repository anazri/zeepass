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
- **File encoding** - Upload files (up to 10MB) and convert to Base64
- **File decoding** - Convert Base64 back to downloadable files
- **Real-time processing** with HTMX
- **Binary data handling** with proper MIME type detection
- **Error handling** for invalid Base64 input and oversized files
- **Clean, intuitive interface**

### 🎨 **User Interface & Experience**
- **Modern Landing Page** with comprehensive feature overview
- **Dark/Light Theme Toggle** with system preference detection
- **Responsive Design** optimized for all device sizes
- **About Us Section** featuring:
  - Our Mission: Making advanced cryptography accessible to everyone
  - Our Expertise: Built by security experts and cryptography specialists  
  - Our Commitment: Continuous innovation and transparency in security
- **Enterprise Contact System** with multiple inquiry types and advanced spam protection
  - Cloud hosting, on-premise installation, technical support inquiries
  - Multi-layered spam protection (honeypot, content filtering, rate limiting)
  - reCAPTCHA v3 integration with score-based verification
  - Professional email notifications via Resend API
- **User Feedback & Survey System** for business intelligence
  - Comprehensive usage analytics and market research
  - Net Promoter Score (NPS) collection
  - Feature request tracking and enterprise interest assessment
- **Smooth Animations** and glassmorphism design effects

## 🏗️ Architecture

### **Backend (Go)**
```
cmd/server/          # Application entry point
internal/
├── handlers/        # HTTP request handlers
│   ├── base64.go    # Base64 encoding/decoding (text & files)
│   ├── chat.go      # Real-time chat WebSocket handlers
│   ├── contact.go   # Enterprise contact system with spam protection
│   ├── encryption.go # Text and file encryption handlers
│   ├── feedback.go  # User survey and feedback collection
│   ├── home.go      # Landing page and core template rendering
│   ├── password.go  # Password generation with strength analysis
│   ├── sshkey.go    # SSH key generation with multiple algorithms
│   ├── survey.go    # Survey page template handler
│   └── view.go      # Encrypted content viewing and validation
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
├── base64.html     # Base64 encoding/decoding tool (text & files)
├── chat-encryption.html  # Real-time encrypted chat
├── file-encryption.html  # File upload and encryption
├── password-generator.html  # Password generation tool
├── ssh-key.html    # SSH key generation tool
├── survey.html     # User feedback and survey collection
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

### **Core Cryptography**
- **AES-256-GCM Encryption**: Military-grade encryption for all data
- **Zero-Knowledge Architecture**: Server never sees unencrypted data
- **SHA-256 PIN Hashing**: Secure PIN protection
- **Automatic Expiration**: Time-based and view-count-based expiry
- **Secure Random Generation**: Cryptographically secure random number generation
- **TLS-Ready**: Designed for HTTPS deployment

### **Advanced Spam & Attack Protection**
- **reCAPTCHA v3 Integration**: Score-based bot detection with intelligent thresholds
- **Multi-layered Spam Protection**:
  - Honeypot field detection for automated form submissions
  - Content-based spam filtering with keyword blacklists
  - Rate limiting (3 requests per 10-minute window per IP)
- **Input Validation**: Comprehensive form validation and sanitization
- **IP Tracking**: Client identification and monitoring for security analysis

## 🏛️ Security Architecture

### **HTMX vs Framework Comparison**

ZeePass uses **HTMX** over client-side frameworks (React/Next.js) for enhanced security:

**✅ HTMX Security Advantages:**
- **Server-Side Cryptography**: All encryption/decryption operations execute server-side in Go
- **Zero Client-Side Crypto**: No JavaScript cryptographic libraries exposed to browser
- **Minimal Attack Surface**: Reduced client-side code minimizes potential vulnerabilities
- **Server-Only Secrets**: Encryption keys never transmitted to or accessible by client
- **XSS Mitigation**: Limited client-side JavaScript reduces XSS-based crypto key extraction risks

**⚠️ Client-Side Framework Risks:**
- JavaScript crypto libraries exposed in browser environment
- Potential crypto keys in client bundles
- Complex dependency chains increase attack surface
- Client-side state management vulnerabilities
- SSR/hydration security considerations

**Security Decision:** HTMX's server-centric approach aligns perfectly with ZeePass's zero-knowledge security model, ensuring all cryptographic operations remain server-side while clients only receive encrypted results.

### **Infrastructure Security**

**Recommended Deployment:** Docker Compose over Kubernetes
- **Simplified Attack Surface**: Fewer moving parts reduce security complexity
- **Container Isolation**: Docker provides process and filesystem isolation
- **Secrets Management**: Docker Compose secrets for encryption keys and Redis passwords
- **Network Security**: Internal container networking isolates services

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

#### **Single Host with Docker Compose**
```bash
# Clone and navigate to the deploy directory
cd deploy

# Start with minimal setup (includes monitoring)
docker-compose -f docker-compose-minimal.yml up -d

# Or full setup with advanced monitoring
docker-compose up -d
```

#### **Multi-Host with Docker Swarm**

**1. Initialize Docker Swarm**
```bash
# On manager node
docker swarm init --advertise-addr <MANAGER-IP>

# On worker nodes (use token from init output)
docker swarm join --token <TOKEN> <MANAGER-IP>:2377
```

**2. Deploy ZeePass Stack**
```bash
# Navigate to deploy directory
cd deploy

# Deploy the stack
docker stack deploy -c docker-swarm-minimal.yml zeepass
```

**3. Manage the Stack**
```bash
# Check stack status
docker stack ps zeepass

# Scale services
docker service scale zeepass_zeepass=3
docker service scale zeepass_nginx=2

# View logs
docker service logs zeepass_zeepass

# Remove stack
docker stack rm zeepass
```

**4. Stack Features**
- **High Availability**: 2 replicas of ZeePass app and Nginx
- **Load Balancing**: Automatic load balancing across replicas  
- **Rolling Updates**: Zero-downtime deployments
- **Health Checks**: Automatic service recovery
- **Monitoring**: GoAccess, Uptime Kuma, and Netdata included

**5. Access Services**
- **ZeePass**: http://your-swarm-ip (load balanced)
- **GoAccess**: http://your-swarm-ip:7890
- **Uptime Kuma**: http://your-swarm-ip:3001  
- **Netdata**: http://your-swarm-ip:19999

#### **Custom Docker Build**
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
- `RESEND_API_KEY`: API key for professional email notifications
- `RECAPTCHA_SECRET_KEY`: reCAPTCHA v3 secret key for spam protection

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
 Software Architect/Engineer & Cybersecurity Expert  
[LinkedIn](https://www.linkedin.com/in/nazri-bin-abdullah-a75b40268/) | [GitHub](https://github.com/anazri)

---

## 🛠️ Tech Stack

- **Backend**: Go 1.24.2
- **Frontend**: HTMX, TailwindCSS, Vanilla JavaScript
- **Database**: Redis
- **Encryption**: AES-256-GCM
- **WebSockets**: Gorilla WebSocket
- **Security**: reCAPTCHA v3, Rate Limiting, Spam Protection
- **Email**: Resend API for professional notifications
- **Development**: Air (live reloading)
- **Deployment**: Docker-ready

## 📊 Project Status

### **Core Features**
- ✅ **Text Encryption** - Complete
- ✅ **File Encryption** - Complete  
- ✅ **Chat Encryption** - Complete
- ✅ **Password Generator** - Complete
- ✅ **SSH Key Generator** - Complete
- ✅ **Base64 Tools** (Text & Files) - Complete
- ✅ **Dark/Light Theme** - Complete
- ✅ **Responsive Design** - Complete

### **Business & Security Features**
- ✅ **Enterprise Contact System** - Complete
- ✅ **User Feedback & Survey System** - Complete
- ✅ **Advanced Spam Protection** - Complete
- ✅ **reCAPTCHA v3 Integration** - Complete
- ✅ **Professional Email Integration** - Complete
- ✅ **Rate Limiting & IP Tracking** - Complete

### **Upcoming Features**

#### **Advanced Crypto Tools**
- 🔄 **PGP/GPG Key Tools** - Generate/import/export OpenPGP keys, encrypt/decrypt/sign messages/files with PGP
- 🔄 **JWT (JSON Web Token) Tools** - Encode/decode/verify JWTs, sign with HS256/RS256/ES256/EdDSA
- 🔄 **Hashing Tools** - Compute SHA-256, SHA-512, BLAKE2, Argon2, MD5 for integrity checks and password hashing
- 🔄 **QR Code Crypto** - Generate QR codes for encrypted messages, passwords, SSH keys with scan/decrypt functionality
- 🔄 **Key Derivation Functions (KDFs)** - PBKDF2, scrypt, Argon2 for secure password-to-key generation
- 🔄 **Digital Signature Tools** - Sign and verify text/files using RSA/ECDSA/Ed25519 for software authenticity
- 🔄 **Certificate & TLS Tools** - Generate CSRs and self-signed X.509 certificates, inspect SSL/TLS certificates
- 🔄 **Mnemonic & Wallet Tools** - Generate BIP39 mnemonics, derive HD wallet keys (BIP32/44), export to ETH/BTC formats
- 🔄 **Steganography Tools** - Hide encrypted text inside images with extraction capabilities
- 🔄 **Entropy & Randomness Tester** - Generate cryptographically secure random numbers with entropy visualization

#### **Platform Features**
- 🔄 **User Authentication** - Planned
- 🔄 **API Endpoints** - Planned
- 🔄 **Mobile App** - Planned
- 🔄 **Analytics Dashboard** - Planned

---

⭐ **Star this repository if you find it useful!**

🐛 **Found a bug?** [Report it here](https://github.com/anazri/zeepass/issues)

💡 **Have a suggestion?** [Let us know!](https://github.com/anazri/zeepass/discussions)