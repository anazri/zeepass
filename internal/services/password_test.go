package services

import (
	"strings"
	"testing"
)

func TestGeneratePasswordRandom(t *testing.T) {
	tests := []struct {
		name string
		opts PasswordOptions
	}{
		{
			"default options",
			PasswordOptions{Length: 12, UseNumbers: true, UseUppercase: true, UseLowercase: true, UseSymbols: true, Type: "random"},
		},
		{
			"numbers only",
			PasswordOptions{Length: 8, UseNumbers: true, Type: "random"},
		},
		{
			"letters only",
			PasswordOptions{Length: 10, UseUppercase: true, UseLowercase: true, Type: "random"},
		},
		{
			"minimum length",
			PasswordOptions{Length: 4, UseNumbers: true, Type: "random"},
		},
		{
			"maximum length",
			PasswordOptions{Length: 64, UseNumbers: true, Type: "random"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GeneratePassword(tt.opts)
			if err != nil {
				t.Fatalf("GeneratePassword failed: %v", err)
			}

			if len(password) != tt.opts.Length {
				t.Errorf("Expected length %d, got %d", tt.opts.Length, len(password))
			}

			if tt.opts.UseNumbers {
				if !containsNumbers(password) {
					t.Error("Password should contain numbers")
				}
			}

			if tt.opts.UseUppercase {
				if !containsUppercase(password) {
					t.Error("Password should contain uppercase letters")
				}
			}

			if tt.opts.UseLowercase {
				if !containsLowercase(password) {
					t.Error("Password should contain lowercase letters")
				}
			}

			if tt.opts.UseSymbols {
				if !containsSymbols(password) {
					t.Error("Password should contain symbols")
				}
			}
		})
	}
}

func TestGeneratePasswordPIN(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"4 digit PIN", 4},
		{"6 digit PIN", 6},
		{"8 digit PIN", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := PasswordOptions{Length: tt.length, Type: "pin"}
			password, err := GeneratePassword(opts)
			if err != nil {
				t.Fatalf("GeneratePassword failed: %v", err)
			}

			if len(password) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(password))
			}

			for _, char := range password {
				if char < '0' || char > '9' {
					t.Errorf("PIN should only contain digits, found %c", char)
				}
			}
		})
	}
}

func TestGeneratePasswordMemorable(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"short memorable", 8},
		{"medium memorable", 16},
		{"long memorable", 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := PasswordOptions{Length: tt.length, Type: "memorable"}
			password, err := GeneratePassword(opts)
			if err != nil {
				t.Fatalf("GeneratePassword failed: %v", err)
			}

			if len(password) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(password))
			}

			hasLetter := false
			hasDigit := false
			for _, char := range password {
				if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
					hasLetter = true
				}
				if char >= '0' && char <= '9' {
					hasDigit = true
				}
			}

			if !hasLetter {
				t.Error("Memorable password should contain letters")
			}
			if !hasDigit {
				t.Error("Memorable password should contain digits")
			}
		})
	}
}

func TestGeneratePasswordEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		opts     PasswordOptions
		expected int
	}{
		{"length too small", PasswordOptions{Length: 1, Type: "random"}, 4},
		{"length too large", PasswordOptions{Length: 100, Type: "random"}, 64},
		{"no charset selected", PasswordOptions{Length: 8, Type: "random"}, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GeneratePassword(tt.opts)
			if err != nil {
				t.Fatalf("GeneratePassword failed: %v", err)
			}

			if len(password) != tt.expected {
				t.Errorf("Expected length %d, got %d", tt.expected, len(password))
			}
		})
	}
}

func TestGeneratePasswordUniqueness(t *testing.T) {
	opts := PasswordOptions{Length: 16, UseNumbers: true, UseUppercase: true, UseLowercase: true, UseSymbols: true, Type: "random"}
	
	passwords := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		password, err := GeneratePassword(opts)
		if err != nil {
			t.Fatalf("GeneratePassword failed: %v", err)
		}

		if passwords[password] {
			t.Errorf("Generated duplicate password: %s", password)
		}
		passwords[password] = true
	}
}

func TestCalculatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{"weak short", "123", "weak"},
		{"weak simple", "password", "weak"},
		{"medium length", "password123", "medium"},
		{"medium mixed", "Password123", "medium"},
		{"strong complete", "Password123!", "strong"},
		{"strong long", "ThisIsAVeryLongPassword123!", "strong"},
		{"empty password", "", "weak"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strength := CalculatePasswordStrength(tt.password)
			if strength != tt.expected {
				t.Errorf("Expected %s, got %s for password %q", tt.expected, strength, tt.password)
			}
		})
	}
}

func TestContainsNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"password123", true},
		{"password", false},
		{"123", true},
		{"", false},
		{"pass1word", true},
	}

	for _, tt := range tests {
		result := containsNumbers(tt.input)
		if result != tt.expected {
			t.Errorf("containsNumbers(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestContainsLowercase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"PASSWORD", false},
		{"Password", true},
		{"password", true},
		{"123", false},
		{"", false},
	}

	for _, tt := range tests {
		result := containsLowercase(tt.input)
		if result != tt.expected {
			t.Errorf("containsLowercase(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestContainsUppercase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"password", false},
		{"Password", true},
		{"PASSWORD", true},
		{"123", false},
		{"", false},
	}

	for _, tt := range tests {
		result := containsUppercase(tt.input)
		if result != tt.expected {
			t.Errorf("containsUppercase(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestContainsSymbols(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"password", false},
		{"password!", true},
		{"pass@word", true},
		{"123", false},
		{"", false},
		{"!@#$%", true},
	}

	for _, tt := range tests {
		result := containsSymbols(tt.input)
		if result != tt.expected {
			t.Errorf("containsSymbols(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestPasswordCharsetValidation(t *testing.T) {
	opts := PasswordOptions{Length: 10, UseNumbers: true, UseUppercase: true, UseLowercase: true, UseSymbols: true, Type: "random"}
	
	password, err := GeneratePassword(opts)
	if err != nil {
		t.Fatalf("GeneratePassword failed: %v", err)
	}

	for _, char := range password {
		validChar := (char >= '0' && char <= '9') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			strings.ContainsRune(Symbols, char)
		
		if !validChar {
			t.Errorf("Invalid character in password: %c", char)
		}
	}
}