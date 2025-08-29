package services

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	testKey := make([]byte, 32)
	rand.Read(testKey)

	tests := []struct {
		name      string
		plaintext string
	}{
		{"simple text", "hello world"},
		{"empty string", ""},
		{"unicode text", "こんにちは世界"},
		{"special chars", "!@#$%^&*()_+-=[]{}|;:,.<>?"},
		{"long text", string(make([]byte, 1000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.plaintext, testKey)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			if encrypted == "" {
				t.Fatal("Encrypted text should not be empty")
			}

			decrypted, err := Decrypt(encrypted, testKey)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("Expected %q, got %q", tt.plaintext, decrypted)
			}
		})
	}
}

func TestEncryptDecryptWithDifferentKeys(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	rand.Read(key1)
	rand.Read(key2)

	plaintext := "secret message"

	encrypted, err := Encrypt(plaintext, key1)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(encrypted, key2)
	if err == nil {
		t.Error("Decrypt should fail with wrong key")
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	testKey := make([]byte, 32)
	rand.Read(testKey)

	tests := []struct {
		name string
		data []byte
	}{
		{"small file", []byte("small file content")},
		{"empty file", []byte{}},
		{"binary data", []byte{0x00, 0x01, 0xFF, 0xFE}},
		{"large file", make([]byte, 10000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := EncryptFile(tt.data, testKey)
			if err != nil {
				t.Fatalf("EncryptFile failed: %v", err)
			}

			if len(encrypted) == 0 && len(tt.data) > 0 {
				t.Fatal("Encrypted data should not be empty for non-empty input")
			}

			decrypted, err := DecryptFile(encrypted, testKey)
			if err != nil {
				t.Fatalf("DecryptFile failed: %v", err)
			}

			if !bytes.Equal(decrypted, tt.data) {
				t.Errorf("Decrypted data doesn't match original")
			}
		})
	}
}

func TestEncryptDecryptFileWithDifferentKeys(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	rand.Read(key1)
	rand.Read(key2)

	data := []byte("secret file content")

	encrypted, err := EncryptFile(data, key1)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	_, err = DecryptFile(encrypted, key2)
	if err == nil {
		t.Error("DecryptFile should fail with wrong key")
	}
}

func TestHashPIN(t *testing.T) {
	tests := []struct {
		name string
		pin  string
	}{
		{"numeric pin", "1234"},
		{"long pin", "123456789"},
		{"empty pin", ""},
		{"alphanumeric", "abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := HashPIN(tt.pin)
			hash2 := HashPIN(tt.pin)

			if hash1 != hash2 {
				t.Error("HashPIN should be deterministic")
			}

			if len(hash1) != 64 {
				t.Errorf("Expected hash length 64, got %d", len(hash1))
			}

			if tt.pin != "" && hash1 == tt.pin {
				t.Error("Hash should be different from original PIN")
			}
		})
	}
}

func TestHashPINDifferentValues(t *testing.T) {
	pin1 := "1234"
	pin2 := "5678"

	hash1 := HashPIN(pin1)
	hash2 := HashPIN(pin2)

	if hash1 == hash2 {
		t.Error("Different PINs should have different hashes")
	}
}

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == id2 {
		t.Error("GenerateID should produce unique IDs")
	}

	if len(id1) != 32 {
		t.Errorf("Expected ID length 32, got %d", len(id1))
	}

	for _, char := range id1 {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("ID should only contain hex characters, found %c", char)
		}
	}
}

func TestDecryptInvalidData(t *testing.T) {
	testKey := make([]byte, 32)
	rand.Read(testKey)

	tests := []struct {
		name       string
		ciphertext string
	}{
		{"invalid base64", "invalid base64 data"},
		{"too short", "dGVzdA=="},
		{"corrupted data", "dGVzdGRhdGF0aGF0aXNsb25nZW5vdWdoYnV0Y29ycnVwdGVk"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.ciphertext, testKey)
			if err == nil {
				t.Error("Decrypt should fail with invalid data")
			}
		})
	}
}

func TestDecryptFileInvalidData(t *testing.T) {
	testKey := make([]byte, 32)
	rand.Read(testKey)

	tests := []struct {
		name string
		data []byte
	}{
		{"too short", []byte("short")},
		{"corrupted data", []byte("this is corrupted ciphertext data that should fail")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptFile(tt.data, testKey)
			if err == nil {
				t.Error("DecryptFile should fail with invalid data")
			}
		})
	}
}