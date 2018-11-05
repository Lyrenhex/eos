package perspectiveapi

// MLRequest is a template for generating Google Perspective API requests
type MLRequest struct {
	Comment         MLComment   `json:"comment"`
	RequestedAttrbs MLAttribute `json:"requestedAttributes"`
	DNS             bool        `json:"doNotStore"`
}

// MLComment : struct, part of MLRequest to serve chat message content
type MLComment struct {
	Text string `json:"text"`
}

// MLAttribute : struct, part of MLRequest to request SEVERE_TOXICITY model results
type MLAttribute struct {
	Attrb MLTOXICITY `json:"SEVERE_TOXICITY"`
}

// MLTOXICITY : struct, part of MLAttribute; empty
type MLTOXICITY struct{}

// MLResponse is a template for parsing Google Perspective API responses
type MLResponse struct {
	AttrbScores AttributeScores `json:"attributeScores"`
}

// AttributeScores : struct, part of MLResponse to parse the SEVERE_TOXICITY results
type AttributeScores struct {
	Toxicity Toxicity `json:"SEVERE_TOXICITY"`
}

// Toxicity : struct, part of Toxicity to parse the summaryScore of the SEVERE_TOXICITY results
type Toxicity struct {
	Summary SummaryScores `json:"summaryScore"`
}

// SummaryScores : struct to parse the summaryScore
type SummaryScores struct {
	Score float64 `json:"value"`
	Type  string  `json:"type"`
}
