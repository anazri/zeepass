package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type SurveyResponse struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	Likelihood     string    `json:"likelihood"`
	Tools          []string  `json:"tools"`
	UseCase        string    `json:"use_case"`
	Concerns       string    `json:"concerns"`
	FeatureRequest string    `json:"feature_request"`
	NPS            int       `json:"nps"`
	Email          string    `json:"email,omitempty"`
	Name           string    `json:"name,omitempty"`
	Updates        bool      `json:"updates"`
	IPAddress      string    `json:"ip_address"`
}

func main() {
	filePath := "data/survey_responses.json"
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Survey data file not found: %s\n", filePath)
		fmt.Println("No survey responses yet.")
		return
	}

	// Read and parse responses
	responses, err := readSurveyResponses(filePath)
	if err != nil {
		log.Fatalf("Error reading survey responses: %v", err)
	}

	if len(responses) == 0 {
		fmt.Println("No survey responses found.")
		return
	}

	// Generate analysis
	fmt.Printf("📊 ZeePass Survey Analysis\n")
	fmt.Printf("=========================\n\n")
	fmt.Printf("Total Responses: %d\n\n", len(responses))

	analyzeLikelihood(responses)
	analyzeTools(responses)
	analyzeUseCases(responses)
	analyzeConcerns(responses)
	analyzeNPS(responses)
	analyzeFeatureRequests(responses)
	analyzeTimeline(responses)

	if hasEmails := countEmailSignups(responses); hasEmails > 0 {
		fmt.Printf("📧 Email Signups: %d (%.1f%%)\n\n", hasEmails, float64(hasEmails)/float64(len(responses))*100)
	}
}

func readSurveyResponses(filePath string) ([]SurveyResponse, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var responses []SurveyResponse
	if err := json.Unmarshal(fileData, &responses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal survey responses: %v", err)
	}

	return responses, nil
}

func analyzeLikelihood(responses []SurveyResponse) {
	fmt.Printf("🎯 Likelihood to Use ZeePass:\n")
	likelihood := make(map[string]int)
	
	for _, r := range responses {
		if r.Likelihood != "" {
			likelihood[r.Likelihood]++
		}
	}

	order := []string{"very_likely", "somewhat_likely", "neutral", "unlikely", "very_unlikely"}
	labels := map[string]string{
		"very_likely":    "Very Likely",
		"somewhat_likely": "Somewhat Likely", 
		"neutral":        "Neutral",
		"unlikely":       "Unlikely",
		"very_unlikely":  "Very Unlikely",
	}

	for _, key := range order {
		if count := likelihood[key]; count > 0 {
			percentage := float64(count) / float64(len(responses)) * 100
			fmt.Printf("  %s: %d (%.1f%%)\n", labels[key], count, percentage)
		}
	}
	fmt.Println()
}

func analyzeTools(responses []SurveyResponse) {
	fmt.Printf("🔧 Most Interesting Tools:\n")
	tools := make(map[string]int)
	
	for _, r := range responses {
		for _, tool := range r.Tools {
			tools[tool]++
		}
	}

	type toolCount struct {
		name  string
		count int
	}
	
	var sortedTools []toolCount
	labels := map[string]string{
		"text_encryption":   "Text Encryption",
		"file_encryption":   "File Encryption",
		"encrypted_chat":    "Encrypted Chat",
		"password_generator": "Password Generator",
		"ssh_key":          "SSH Key Generator",
		"base64":           "Base64 Tools",
	}

	for tool, count := range tools {
		sortedTools = append(sortedTools, toolCount{tool, count})
	}

	sort.Slice(sortedTools, func(i, j int) bool {
		return sortedTools[i].count > sortedTools[j].count
	})

	for _, tc := range sortedTools {
		percentage := float64(tc.count) / float64(len(responses)) * 100
		label := labels[tc.name]
		if label == "" {
			label = tc.name
		}
		fmt.Printf("  %s: %d (%.1f%%)\n", label, tc.count, percentage)
	}
	fmt.Println()
}

func analyzeUseCases(responses []SurveyResponse) {
	fmt.Printf("💡 Primary Use Cases:\n")
	useCases := make(map[string]int)
	
	for _, r := range responses {
		if r.UseCase != "" {
			useCases[r.UseCase]++
		}
	}

	labels := map[string]string{
		"personal_privacy":   "Personal Privacy & Security",
		"team_collaboration": "Team/Work Collaboration",
		"development_it":     "Development & IT Tasks",
		"sensitive_sharing":  "Sharing Sensitive Information",
		"learning":          "Learning About Encryption",
		"other":             "Other",
	}

	type useCaseCount struct {
		name  string
		count int
	}
	
	var sortedUseCases []useCaseCount
	for useCase, count := range useCases {
		sortedUseCases = append(sortedUseCases, useCaseCount{useCase, count})
	}

	sort.Slice(sortedUseCases, func(i, j int) bool {
		return sortedUseCases[i].count > sortedUseCases[j].count
	})

	for _, uc := range sortedUseCases {
		percentage := float64(uc.count) / float64(len(responses)) * 100
		label := labels[uc.name]
		if label == "" {
			label = uc.name
		}
		fmt.Printf("  %s: %d (%.1f%%)\n", label, uc.count, percentage)
	}
	fmt.Println()
}

func analyzeConcerns(responses []SurveyResponse) {
	fmt.Printf("⚠️  Main Concerns:\n")
	concerns := make(map[string]int)
	
	for _, r := range responses {
		if r.Concerns != "" {
			concerns[r.Concerns]++
		}
	}

	labels := map[string]string{
		"data_privacy":        "Data Privacy & Trust",
		"ease_of_use":        "Ease of Use",
		"feature_completeness": "Feature Completeness",
		"reliability":        "Reliability & Uptime",
		"cost":               "Cost",
		"other":              "Other",
	}

	type concernCount struct {
		name  string
		count int
	}
	
	var sortedConcerns []concernCount
	for concern, count := range concerns {
		sortedConcerns = append(sortedConcerns, concernCount{concern, count})
	}

	sort.Slice(sortedConcerns, func(i, j int) bool {
		return sortedConcerns[i].count > sortedConcerns[j].count
	})

	for _, cc := range sortedConcerns {
		percentage := float64(cc.count) / float64(len(responses)) * 100
		label := labels[cc.name]
		if label == "" {
			label = cc.name
		}
		fmt.Printf("  %s: %d (%.1f%%)\n", label, cc.count, percentage)
	}
	fmt.Println()
}

func analyzeNPS(responses []SurveyResponse) {
	fmt.Printf("📊 Net Promoter Score (NPS):\n")
	
	npsScores := make([]int, 0, len(responses))
	npsCount := make(map[int]int)
	
	for _, r := range responses {
		if r.NPS >= 0 && r.NPS <= 10 {
			npsScores = append(npsScores, r.NPS)
			npsCount[r.NPS]++
		}
	}

	if len(npsScores) == 0 {
		fmt.Printf("  No NPS data available\n\n")
		return
	}

	// Calculate NPS
	promoters := 0
	passives := 0
	detractors := 0

	for _, score := range npsScores {
		if score >= 9 {
			promoters++
		} else if score >= 7 {
			passives++
		} else {
			detractors++
		}
	}

	nps := float64(promoters-detractors) / float64(len(npsScores)) * 100

	fmt.Printf("  Total NPS Responses: %d\n", len(npsScores))
	fmt.Printf("  Promoters (9-10): %d (%.1f%%)\n", promoters, float64(promoters)/float64(len(npsScores))*100)
	fmt.Printf("  Passives (7-8): %d (%.1f%%)\n", passives, float64(passives)/float64(len(npsScores))*100)
	fmt.Printf("  Detractors (0-6): %d (%.1f%%)\n", detractors, float64(detractors)/float64(len(npsScores))*100)
	fmt.Printf("  📈 NPS Score: %.1f\n", nps)

	// NPS interpretation
	if nps >= 70 {
		fmt.Printf("  🎉 Excellent! (World-class NPS)\n")
	} else if nps >= 50 {
		fmt.Printf("  ✅ Great! (Excellent NPS)\n")
	} else if nps >= 30 {
		fmt.Printf("  👍 Good (Above average NPS)\n")
	} else if nps >= 0 {
		fmt.Printf("  ⚠️  Needs improvement (Below average NPS)\n")
	} else {
		fmt.Printf("  ❌ Critical - immediate action needed\n")
	}

	fmt.Println()
}

func analyzeFeatureRequests(responses []SurveyResponse) {
	fmt.Printf("💭 Feature Requests:\n")
	
	requests := make([]string, 0)
	for _, r := range responses {
		if strings.TrimSpace(r.FeatureRequest) != "" {
			requests = append(requests, fmt.Sprintf("• %s", strings.TrimSpace(r.FeatureRequest)))
		}
	}

	if len(requests) == 0 {
		fmt.Printf("  No feature requests yet\n\n")
		return
	}

	fmt.Printf("  Total Requests: %d\n", len(requests))
	for _, req := range requests {
		fmt.Printf("  %s\n", req)
	}
	fmt.Println()
}

func analyzeTimeline(responses []SurveyResponse) {
	fmt.Printf("📅 Response Timeline:\n")
	
	if len(responses) == 0 {
		return
	}

	// Sort by timestamp
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Timestamp.Before(responses[j].Timestamp)
	})

	first := responses[0].Timestamp
	last := responses[len(responses)-1].Timestamp

	fmt.Printf("  First Response: %s\n", first.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Latest Response: %s\n", last.Format("2006-01-02 15:04:05"))
	
	if len(responses) > 1 {
		duration := last.Sub(first)
		fmt.Printf("  Time Span: %s\n", duration.String())
	}
	fmt.Println()
}

func countEmailSignups(responses []SurveyResponse) int {
	count := 0
	for _, r := range responses {
		if strings.TrimSpace(r.Email) != "" {
			count++
		}
	}
	return count
}