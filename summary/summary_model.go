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

	WhatYouAgreeOn string `json:"what_you_agree_on"`
}

func NewSummary(metadata SummaryMetadata, whatYouAgreeOn string) *Summary {
	return &Summary{SummaryMetadata: metadata, WhatYouAgreeOn: whatYouAgreeOn}
}
