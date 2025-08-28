package services

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGenerateSSHKeyRSA(t *testing.T) {
	tests := []struct {
		name       string
		bits       int
		passphrase string
		comment    string
	}{
		{"RSA 2048", 2048, "", "test@example.com"},
		{"RSA 3072", 3072, "", "test"},
		{"RSA 4096", 4096, "", ""},
		{"RSA with passphrase", 2048, "testpass", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := SSHKeyOptions{
				Type:       "rsa",
				Length:     tt.bits,
				Passphrase: tt.passphrase,
				Comment:    tt.comment,
			}

			keyPair, err := GenerateSSHKey(opts)
			if err != nil {
				t.Fatalf("GenerateSSHKey failed: %v", err)
			}

			if keyPair.PrivateKey == "" {
				t.Error("Private key should not be empty")
			}

			if keyPair.PublicKey == "" {
				t.Error("Public key should not be empty")
			}

			if !strings.HasPrefix(keyPair.PublicKey, "ssh-rsa") {
				t.Error("Public key should start with 'ssh-rsa'")
			}

			expectedComment := tt.comment
			if expectedComment == "" {
				expectedComment = "noname"
			}
			if !strings.HasSuffix(keyPair.PublicKey, " "+expectedComment) {
				t.Errorf("Public key should end with comment: %s", expectedComment)
			}

			block, _ := pem.Decode([]byte(keyPair.PrivateKey))
			if block == nil {
				t.Fatal("Failed to decode PEM block")
			}

			if tt.passphrase != "" {
				if !x509.IsEncryptedPEMBlock(block) {
					t.Error("Private key should be encrypted when passphrase provided")
				}
			} else {
				if x509.IsEncryptedPEMBlock(block) {
					t.Error("Private key should not be encrypted when no passphrase provided")
				}
			}
		})
	}
}

func TestGenerateSSHKeyEd25519(t *testing.T) {
	tests := []struct {
		name       string
		passphrase string
		comment    string
	}{
		{"Ed25519 basic", "", "test@example.com"},
		{"Ed25519 with passphrase", "testpass", "test"},
		{"Ed25519 no comment", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := SSHKeyOptions{
				Type:       "ed25519",
				Passphrase: tt.passphrase,
				Comment:    tt.comment,
			}

			keyPair, err := GenerateSSHKey(opts)
			if err != nil {
				t.Fatalf("GenerateSSHKey failed: %v", err)
			}

			if keyPair.PrivateKey == "" {
				t.Error("Private key should not be empty")
			}

			if keyPair.PublicKey == "" {
				t.Error("Public key should not be empty")
			}

			if !strings.HasPrefix(keyPair.PublicKey, "ssh-ed25519") {
				t.Error("Public key should start with 'ssh-ed25519'")
			}

			expectedComment := tt.comment
			if expectedComment == "" {
				expectedComment = "noname"
			}
			if !strings.HasSuffix(keyPair.PublicKey, " "+expectedComment) {
				t.Errorf("Public key should end with comment: %s", expectedComment)
			}

			block, _ := pem.Decode([]byte(keyPair.PrivateKey))
			if block == nil {
				t.Fatal("Failed to decode PEM block")
			}

			if tt.passphrase != "" {
				if !x509.IsEncryptedPEMBlock(block) {
					t.Error("Private key should be encrypted when passphrase provided")
				}
			}
		})
	}
}

func TestGenerateSSHKeyECDSA(t *testing.T) {
	tests := []struct {
		name       string
		bits       int
		passphrase string
		comment    string
	}{
		{"ECDSA 256", 256, "", "test@example.com"},
		{"ECDSA 384", 384, "", "test"},
		{"ECDSA 521", 521, "", ""},
		{"ECDSA with passphrase", 256, "testpass", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := SSHKeyOptions{
				Type:       "ecdsa",
				Length:     tt.bits,
				Passphrase: tt.passphrase,
				Comment:    tt.comment,
			}

			keyPair, err := GenerateSSHKey(opts)
			if err != nil {
				t.Fatalf("GenerateSSHKey failed: %v", err)
			}

			if keyPair.PrivateKey == "" {
				t.Error("Private key should not be empty")
			}

			if keyPair.PublicKey == "" {
				t.Error("Public key should not be empty")
			}

			expectedPrefix := "ecdsa-sha2-"
			if !strings.HasPrefix(keyPair.PublicKey, expectedPrefix) {
				t.Errorf("Public key should start with '%s'", expectedPrefix)
			}

			expectedComment := tt.comment
			if expectedComment == "" {
				expectedComment = "noname"
			}
			if !strings.HasSuffix(keyPair.PublicKey, " "+expectedComment) {
				t.Errorf("Public key should end with comment: %s", expectedComment)
			}

			block, _ := pem.Decode([]byte(keyPair.PrivateKey))
			if block == nil {
				t.Fatal("Failed to decode PEM block")
			}

			if tt.passphrase != "" {
				if !x509.IsEncryptedPEMBlock(block) {
					t.Error("Private key should be encrypted when passphrase provided")
				}
			}
		})
	}
}

func TestGenerateSSHKeyInvalidType(t *testing.T) {
	opts := SSHKeyOptions{
		Type:   "invalid",
		Length: 2048,
	}

	_, err := GenerateSSHKey(opts)
	if err == nil {
		t.Error("GenerateSSHKey should fail with invalid key type")
	}

	if !strings.Contains(err.Error(), "unsupported key type") {
		t.Errorf("Error should mention unsupported key type, got: %v", err)
	}
}

func TestValidateSSHKeyOptions(t *testing.T) {
	tests := []struct {
		name      string
		opts      SSHKeyOptions
		expectErr bool
	}{
		{"valid RSA 2048", SSHKeyOptions{Type: "rsa", Length: 2048}, false},
		{"valid RSA 3072", SSHKeyOptions{Type: "rsa", Length: 3072}, false},
		{"valid RSA 4096", SSHKeyOptions{Type: "rsa", Length: 4096}, false},
		{"invalid RSA too small", SSHKeyOptions{Type: "rsa", Length: 1024}, true},
		{"invalid RSA too large", SSHKeyOptions{Type: "rsa", Length: 8192}, true},
		{"invalid RSA odd size", SSHKeyOptions{Type: "rsa", Length: 2500}, true},
		{"valid Ed25519", SSHKeyOptions{Type: "ed25519", Length: 0}, false},
		{"valid ECDSA 256", SSHKeyOptions{Type: "ecdsa", Length: 256}, false},
		{"valid ECDSA 384", SSHKeyOptions{Type: "ecdsa", Length: 384}, false},
		{"valid ECDSA 521", SSHKeyOptions{Type: "ecdsa", Length: 521}, false},
		{"invalid ECDSA", SSHKeyOptions{Type: "ecdsa", Length: 512}, true},
		{"invalid type", SSHKeyOptions{Type: "invalid", Length: 2048}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSSHKeyOptions(tt.opts)
			if tt.expectErr && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestSSHKeyPairValidity(t *testing.T) {
	tests := []struct {
		name string
		opts SSHKeyOptions
	}{
		{"RSA key", SSHKeyOptions{Type: "rsa", Length: 2048, Comment: "test"}},
		{"Ed25519 key", SSHKeyOptions{Type: "ed25519", Comment: "test"}},
		{"ECDSA key", SSHKeyOptions{Type: "ecdsa", Length: 256, Comment: "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyPair, err := GenerateSSHKey(tt.opts)
			if err != nil {
				t.Fatalf("GenerateSSHKey failed: %v", err)
			}

			publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(keyPair.PublicKey))
			if err != nil {
				t.Fatalf("Failed to parse public key: %v", err)
			}

			if publicKey == nil {
				t.Error("Public key should be parseable")
			}

			block, _ := pem.Decode([]byte(keyPair.PrivateKey))
			if block == nil {
				t.Fatal("Private key should be valid PEM")
			}
		})
	}
}

func TestGenerateSSHKeyUniqueness(t *testing.T) {
	opts := SSHKeyOptions{Type: "ed25519", Comment: "test"}
	
	keyPair1, err := GenerateSSHKey(opts)
	if err != nil {
		t.Fatalf("GenerateSSHKey failed: %v", err)
	}

	keyPair2, err := GenerateSSHKey(opts)
	if err != nil {
		t.Fatalf("GenerateSSHKey failed: %v", err)
	}

	if keyPair1.PrivateKey == keyPair2.PrivateKey {
		t.Error("Generated keys should be unique")
	}

	if keyPair1.PublicKey == keyPair2.PublicKey {
		t.Error("Generated public keys should be unique")
	}
}

func TestECDSAInvalidKeyLength(t *testing.T) {
	opts := SSHKeyOptions{
		Type:   "ecdsa",
		Length: 999,
	}

	_, err := GenerateSSHKey(opts)
	if err == nil {
		t.Error("Should fail with invalid ECDSA key length")
	}

	if !strings.Contains(err.Error(), "unsupported ECDSA key length") {
		t.Errorf("Error should mention unsupported ECDSA key length, got: %v", err)
	}
}