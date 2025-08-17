package services

import (
	"crypto/rand"
	"math/big"
	"strings"
)

type PasswordOptions struct {
	Length       int  `json:"length"`
	UseNumbers   bool `json:"use_numbers"`
	UseUppercase bool `json:"use_uppercase"`
	UseLowercase bool `json:"use_lowercase"`
	UseSymbols   bool `json:"use_symbols"`
	Type         string `json:"type"` // "random", "memorable", "pin"
}

const (
	Numbers    = "0123456789"
	Uppercase  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase  = "abcdefghijklmnopqrstuvwxyz"
	Symbols    = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

var memorableWords = []string{
	"cat", "dog", "sun", "moon", "star", "tree", "rock", "fire", "wind", "sea",
	"blue", "red", "gold", "fast", "slow", "big", "small", "light", "dark", "bright",
	"apple", "book", "car", "door", "fish", "game", "house", "ice", "key", "lamp",
}

func GeneratePassword(opts PasswordOptions) (string, error) {
	if opts.Length < 4 {
		opts.Length = 4
	}
	if opts.Length > 64 {
		opts.Length = 64
	}

	switch opts.Type {
	case "pin":
		return generatePIN(opts.Length)
	case "memorable":
		return generateMemorablePassword(opts.Length)
	default:
		return generateRandomPassword(opts)
	}
}

func generateRandomPassword(opts PasswordOptions) (string, error) {
	charset := ""
	
	if opts.UseNumbers {
		charset += Numbers
	}
	if opts.UseUppercase {
		charset += Uppercase
	}
	if opts.UseLowercase {
		charset += Lowercase
	}
	if opts.UseSymbols {
		charset += Symbols
	}
	
	// Default to numbers if no charset selected
	if charset == "" {
		charset = Numbers
	}
	
	password := make([]byte, opts.Length)
	charsetLen := big.NewInt(int64(len(charset)))
	
	for i := 0; i < opts.Length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}
	
	return string(password), nil
}

func generatePIN(length int) (string, error) {
	pin := make([]byte, length)
	charsetLen := big.NewInt(int64(len(Numbers)))
	
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		pin[i] = Numbers[randomIndex.Int64()]
	}
	
	return string(pin), nil
}

func generateMemorablePassword(length int) (string, error) {
	var password strings.Builder
	remaining := length
	
	// Add words until we approach the desired length
	for remaining > 4 {
		// Select random word
		wordIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(memorableWords))))
		if err != nil {
			return "", err
		}
		
		word := memorableWords[wordIndex.Int64()]
		
		// Capitalize first letter
		capitalizedWord := strings.ToUpper(string(word[0])) + word[1:]
		
		if password.Len()+len(capitalizedWord) <= length {
			password.WriteString(capitalizedWord)
			remaining -= len(capitalizedWord)
		} else {
			break
		}
	}
	
	// Fill remaining length with numbers
	for remaining > 0 {
		digitIndex, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		password.WriteByte('0' + byte(digitIndex.Int64()))
		remaining--
	}
	
	result := password.String()
	
	// If password is too long, trim it
	if len(result) > length {
		result = result[:length]
	}
	
	// If password is too short, pad with numbers
	for len(result) < length {
		digitIndex, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result += string('0' + byte(digitIndex.Int64()))
	}
	
	return result, nil
}

func CalculatePasswordStrength(password string) string {
	score := 0
	
	// Length score
	if len(password) >= 12 {
		score += 2
	} else if len(password) >= 8 {
		score += 1
	}
	
	// Character variety score
	if containsNumbers(password) {
		score += 1
	}
	if containsLowercase(password) {
		score += 1
	}
	if containsUppercase(password) {
		score += 1
	}
	if containsSymbols(password) {
		score += 1
	}
	
	if score >= 5 {
		return "strong"
	} else if score >= 3 {
		return "medium"
	} else {
		return "weak"
	}
}

func containsNumbers(s string) bool {
	for _, char := range s {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, char := range s {
		if char >= 'a' && char <= 'z' {
			return true
		}
	}
	return false
}

func containsUppercase(s string) bool {
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

func containsSymbols(s string) bool {
	for _, char := range s {
		if strings.ContainsRune(Symbols, char) {
			return true
		}
	}
	return false
}