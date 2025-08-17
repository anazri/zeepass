package services

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
)

type SSHKeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type SSHKeyOptions struct {
	Type       string `json:"type"`       // "rsa", "ed25519", "ecdsa"
	Length     int    `json:"length"`     // Key length in bits
	Passphrase string `json:"passphrase"` // Optional passphrase
	Comment    string `json:"comment"`    // Comment for the public key
}

// GenerateSSHKey generates an SSH key pair based on the provided options
func GenerateSSHKey(opts SSHKeyOptions) (*SSHKeyPair, error) {
	// Set default comment if empty
	if opts.Comment == "" {
		opts.Comment = "noname"
	}

	switch opts.Type {
	case "rsa":
		return generateRSAKey(opts.Length, opts.Passphrase, opts.Comment)
	case "ed25519":
		return generateEd25519Key(opts.Passphrase, opts.Comment)
	case "ecdsa":
		return generateECDSAKey(opts.Length, opts.Passphrase, opts.Comment)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", opts.Type)
	}
}

// generateRSAKey generates an RSA SSH key pair
func generateRSAKey(bits int, passphrase string, comment string) (*SSHKeyPair, error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Encrypt private key with passphrase if provided
	if passphrase != "" {
		privateKeyPEM, err = x509.EncryptPEMBlock(rand.Reader, privateKeyPEM.Type, privateKeyPEM.Bytes, []byte(passphrase), x509.PEMCipherAES256)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %v", err)
		}
	}

	privateKeyStr := string(pem.EncodeToMemory(privateKeyPEM))

	// Generate SSH public key
	sshPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SSH public key: %v", err)
	}

	publicKeyStr := fmt.Sprintf("%s %s", sshPublicKey.Type(), string(ssh.MarshalAuthorizedKey(sshPublicKey)))
	publicKeyStr = publicKeyStr[:len(publicKeyStr)-1] + " " + comment // Replace newline with comment

	return &SSHKeyPair{
		PrivateKey: privateKeyStr,
		PublicKey:  publicKeyStr,
	}, nil
}

// generateEd25519Key generates an Ed25519 SSH key pair
func generateEd25519Key(passphrase string, comment string) (*SSHKeyPair, error) {
	// Generate Ed25519 private key
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 key: %v", err)
	}

	// Encode private key to PEM format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Encrypt private key with passphrase if provided
	if passphrase != "" {
		privateKeyPEM, err = x509.EncryptPEMBlock(rand.Reader, privateKeyPEM.Type, privateKeyPEM.Bytes, []byte(passphrase), x509.PEMCipherAES256)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %v", err)
		}
	}

	privateKeyStr := string(pem.EncodeToMemory(privateKeyPEM))

	// Generate SSH public key
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SSH public key: %v", err)
	}

	publicKeyStr := fmt.Sprintf("%s %s", sshPublicKey.Type(), string(ssh.MarshalAuthorizedKey(sshPublicKey)))
	publicKeyStr = publicKeyStr[:len(publicKeyStr)-1] + " " + comment // Replace newline with comment

	return &SSHKeyPair{
		PrivateKey: privateKeyStr,
		PublicKey:  publicKeyStr,
	}, nil
}

// generateECDSAKey generates an ECDSA SSH key pair
func generateECDSAKey(bits int, passphrase string, comment string) (*SSHKeyPair, error) {
	var curve elliptic.Curve

	// Select curve based on key length
	switch bits {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("unsupported ECDSA key length: %d (supported: 256, 384, 521)", bits)
	}

	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key: %v", err)
	}

	// Encode private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Encrypt private key with passphrase if provided
	if passphrase != "" {
		privateKeyPEM, err = x509.EncryptPEMBlock(rand.Reader, privateKeyPEM.Type, privateKeyPEM.Bytes, []byte(passphrase), x509.PEMCipherAES256)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %v", err)
		}
	}

	privateKeyStr := string(pem.EncodeToMemory(privateKeyPEM))

	// Generate SSH public key
	sshPublicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SSH public key: %v", err)
	}

	publicKeyStr := fmt.Sprintf("%s %s", sshPublicKey.Type(), string(ssh.MarshalAuthorizedKey(sshPublicKey)))
	publicKeyStr = publicKeyStr[:len(publicKeyStr)-1] + " " + comment // Replace newline with comment

	return &SSHKeyPair{
		PrivateKey: privateKeyStr,
		PublicKey:  publicKeyStr,
	}, nil
}

// ValidateSSHKeyOptions validates the SSH key generation options
func ValidateSSHKeyOptions(opts SSHKeyOptions) error {
	switch opts.Type {
	case "rsa":
		if opts.Length < 2048 || opts.Length > 4096 {
			return fmt.Errorf("RSA key length must be between 2048 and 4096 bits")
		}
		if opts.Length != 2048 && opts.Length != 3072 && opts.Length != 4096 {
			return fmt.Errorf("RSA key length must be 2048, 3072, or 4096 bits")
		}
	case "ed25519":
		// Ed25519 keys have a fixed length
		opts.Length = 256
	case "ecdsa":
		if opts.Length != 256 && opts.Length != 384 && opts.Length != 521 {
			return fmt.Errorf("ECDSA key length must be 256, 384, or 521 bits")
		}
	default:
		return fmt.Errorf("unsupported key type: %s (supported: rsa, ed25519, ecdsa)", opts.Type)
	}

	return nil
}
