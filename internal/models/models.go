package models

// IPO represents an Initial Public Offering
type IPO struct {
	Name                string
	OpenDate            string
	CloseDate           string
	PriceRange          string
	GMP                 string // Grey Market Premium
	SubscriptionDetails string
	LotSize             string
	Exchange            string // NSE, BSE, or both
	CompanyInfo         string
}

