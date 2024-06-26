package summary

type SummaryMetadata struct {
	ServiceProvider string `json:"service_provider"`
	EffectiveDate   string `json:"effective_date"`
}

func NewSummaryMetadata(serviceProvider, effectiveDate string) *SummaryMetadata {
	return &SummaryMetadata{ServiceProvider: serviceProvider, EffectiveDate: effectiveDate}
}

type Summary struct {
	SummaryMetadata

	Content string `json:"content"`
}

func NewSummary(metadata SummaryMetadata, content string) *Summary {
	return &Summary{
		SummaryMetadata: metadata,
		Content:         content,
	}
}
