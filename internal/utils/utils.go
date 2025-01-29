package utils

type AISummaryResponse struct {
	Summary string `json:"summary"`
}

// ContentPart represents a part of the content.
type ContentPart struct {
	Text string `json:"text"`
}

// Content represents the content structure.
type Content struct {
	Parts []ContentPart `json:"parts"`
}

// GenerateContentRequest represents the API request structure.
type GenerateContentRequest struct {
	Contents []Content `json:"contents"`
}

// GenerateContentResponse represents the API response structure.
type GenerateContentResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason     string `json:"finishReason"`
		CitationMetadata struct {
			CitationSources []struct {
				StartIndex int    `json:"startIndex"`
				EndIndex   int    `json:"endIndex"`
				URI        string `json:"uri"`
			} `json:"citationSources"`
		} `json:"citationMetadata"`
		AvgLogprobs float64 `json:"avgLogprobs"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
	ModelVersion string `json:"modelVersion"`
}
