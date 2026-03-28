package aiops

type analyzeRequest struct {
	SourceType string `json:"sourceType"`
	SourceID   uint64 `json:"sourceId"`
	Text       string `json:"text"`
}

type termScore struct {
	Term     string `json:"term"`
	Severity int    `json:"severity"`
	Category string `json:"category"`
	Hits     int    `json:"hits"`
}
