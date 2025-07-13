package models

type AnalysisResult struct {
	ID                 uint `gorm:"primaryKey"`
	URLID              uint `gorm:"uniqueIndex"` // Ensures only one analysis per URL
	HTMLVersion        string
	Title              string
	H1Count            int
	H2Count            int
	H3Count            int
	H4Count            int
	H5Count            int
	H6Count            int
	InternalLinksCount int
	ExternalLinksCount int
	BrokenLinksCount   int
	HasLoginForm       bool
}
