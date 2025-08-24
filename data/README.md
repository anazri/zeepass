# Survey Data

This directory contains survey response data from the ZeePass feedback form.

## Files

- `survey_responses.json` - Survey responses in JSON array format

## Data Format

The file contains an array of survey response objects:

```json
[
  {
    "id": "survey_1642123456_789",
    "timestamp": "2025-01-14T10:30:45Z",
    "likelihood": "very_likely",
    "tools": ["text_encryption", "file_encryption"],
    "use_case": "personal_privacy",
    "concerns": "data_privacy",
    "feature_request": "Mobile app support",
    "nps": 9,
    "email": "user@example.com",
    "name": "John Doe",
    "updates": true,
    "ip_address": "192.168.1.100"
  },
  {
    "id": "survey_1642123460_123",
    "timestamp": "2025-01-14T11:15:30Z",
    "likelihood": "somewhat_likely",
    "tools": ["password_generator", "ssh_key"],
    "use_case": "development_it",
    "concerns": "ease_of_use",
    "feature_request": "Better documentation",
    "nps": 7,
    "email": "",
    "name": "",
    "updates": false,
    "ip_address": "192.168.1.101"
  }
]
```

## Analysis

To analyze survey data, run:

```bash
# From the project root
go run scripts/analyze_survey.go
```

This will generate a comprehensive analysis including:
- Response distribution by likelihood to use
- Most popular tools
- Primary use cases
- Main user concerns
- Net Promoter Score (NPS) analysis
- Feature request summary
- Response timeline

## Data Privacy

- IP addresses are collected for basic analytics but not used for tracking
- Email addresses are only stored if users opt-in
- No personal data is shared with third parties
- All data is stored locally on your server

## Manual Analysis

You can also analyze the data manually using `jq`:

```bash
# Count total responses
jq 'length' data/survey_responses.json

# View latest response
jq '.[-1]' data/survey_responses.json

# Extract all feature requests
jq -r '.[].feature_request | select(. != "")' data/survey_responses.json

# Count NPS scores
jq '.[].nps' data/survey_responses.json | sort -n | uniq -c

# Get all tools mentioned
jq -r '.[].tools[]' data/survey_responses.json | sort | uniq -c

# Filter by likelihood
jq '.[] | select(.likelihood == "very_likely")' data/survey_responses.json
```