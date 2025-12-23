package fetcher

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"stock-picker/internal/models"
)

// IPOFetcher fetches IPO data from various sources
type IPOFetcher struct {
	client *http.Client
}

// NewIPOFetcher creates a new IPOFetcher instance
func NewIPOFetcher() *IPOFetcher {
	return &IPOFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchOpenIPOs fetches IPO data from multiple sources
func (f *IPOFetcher) FetchOpenIPOs() ([]models.IPO, error) {
	var allIPOs []models.IPO
	iposMap := make(map[string]*models.IPO)

	// Fetch from multiple sources and merge
	sources := []func() ([]models.IPO, error){
		f.fetchFromChittorgarh,
		f.fetchFromIPOAlert,
	}

	for _, source := range sources {
		ipos, err := source()
		if err != nil {
			log.Printf("Error fetching from source: %v", err)
			continue
		}
		for _, ipo := range ipos {
			// Merge IPOs by name (case-insensitive)
			key := strings.ToLower(ipo.Name)
			if existing, exists := iposMap[key]; exists {
				// Merge data, prefer non-empty values
				if ipo.GMP != "" && existing.GMP == "" {
					existing.GMP = ipo.GMP
				}
				if ipo.SubscriptionDetails != "" && existing.SubscriptionDetails == "" {
					existing.SubscriptionDetails = ipo.SubscriptionDetails
				}
				if ipo.LotSize != "" && existing.LotSize == "" {
					existing.LotSize = ipo.LotSize
				}
			} else {
				iposMap[key] = &ipo
			}
		}
	}

	// Convert map to slice
	for _, ipo := range iposMap {
		allIPOs = append(allIPOs, *ipo)
	}

	// Filter only open IPOs
	openIPOs := f.filterOpenIPOs(allIPOs)

	return openIPOs, nil
}

// fetchFromChittorgarh fetches IPO data from Chittorgarh.com
func (f *IPOFetcher) fetchFromChittorgarh() ([]models.IPO, error) {
	url := "https://www.chittorgarh.com/report/ipo_list_main/ipo_list_main.asp"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Chittorgarh: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var ipos []models.IPO
	// Try multiple table selectors as website structure may vary
	selectors := []string{
		"table tbody tr",
		"table tr",
		".ipo-table tr",
		"#ipo-table tr",
	}

	var found bool
	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			// Skip header rows
			if i == 0 && strings.Contains(strings.ToLower(s.Text()), "company") {
				return
			}

			ipo := models.IPO{}
			cells := s.Find("td")

			if cells.Length() < 3 {
				return
			}

			// Extract IPO name (usually first column)
			name := strings.TrimSpace(cells.Eq(0).Text())
			if name == "" || strings.EqualFold(name, "company name") {
				return
			}
			ipo.Name = name

			// Extract dates (try different formats)
			dates := strings.TrimSpace(cells.Eq(1).Text())
			if dates != "" {
				// Try "to" separator
				if parts := strings.Split(dates, " to "); len(parts) == 2 {
					ipo.OpenDate = strings.TrimSpace(parts[0])
					ipo.CloseDate = strings.TrimSpace(parts[1])
				} else if parts := strings.Split(dates, "-"); len(parts) == 2 {
					ipo.OpenDate = strings.TrimSpace(parts[0])
					ipo.CloseDate = strings.TrimSpace(parts[1])
				}
			}

			// Extract price range
			if cells.Length() > 2 {
				ipo.PriceRange = strings.TrimSpace(cells.Eq(2).Text())
			}

			// Extract lot size
			if cells.Length() > 3 {
				ipo.LotSize = strings.TrimSpace(cells.Eq(3).Text())
			}

			// Extract exchange
			if cells.Length() > 4 {
				ipo.Exchange = strings.TrimSpace(cells.Eq(4).Text())
			}

			// Extract subscription details
			if cells.Length() > 5 {
				ipo.SubscriptionDetails = strings.TrimSpace(cells.Eq(5).Text())
			}

			if ipo.Name != "" {
				ipos = append(ipos, ipo)
				found = true
			}
		})

		if found {
			break
		}
	}

	if len(ipos) == 0 {
		return nil, fmt.Errorf("no IPO data found on page")
	}

	return ipos, nil
}

// fetchFromIPOAlert fetches IPO data from IPOAlert.in (alternative source)
func (f *IPOFetcher) fetchFromIPOAlert() ([]models.IPO, error) {
	// This is a placeholder - you may need to adjust based on actual API or scraping
	// For now, we'll try to fetch GMP data from a GMP-specific source
	return f.fetchGMPData()
}

// fetchGMPData fetches Grey Market Premium data
func (f *IPOFetcher) fetchGMPData() ([]models.IPO, error) {
	// Try to fetch from GMP share or similar sources
	// This is a simplified version - you may need to adjust based on actual website structure
	url := "https://www.gmpshare.com/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := f.client.Do(req)
	if err != nil {
		// If this fails, return empty slice (not critical)
		return []models.IPO{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []models.IPO{}, nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []models.IPO{}, nil
	}

	var ipos []models.IPO
	// This is a placeholder - adjust selectors based on actual website structure
	doc.Find("table tr, .ipo-item, .gmp-item").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find("td:first-child, .ipo-name, .company-name").First().Text())
		if name == "" {
			return
		}

		gmp := strings.TrimSpace(s.Find("td:nth-child(2), .gmp-value").First().Text())

		// Try to find existing IPO or create new one
		ipo := models.IPO{
			Name: name,
			GMP:  gmp,
		}
		ipos = append(ipos, ipo)
	})

	return ipos, nil
}

// filterOpenIPOs filters IPOs that are currently open
func (f *IPOFetcher) filterOpenIPOs(ipos []models.IPO) []models.IPO {
	now := time.Now()
	var openIPOs []models.IPO

	for _, ipo := range ipos {
		if f.isIPOOpen(ipo, now) {
			openIPOs = append(openIPOs, ipo)
		}
	}

	return openIPOs
}

// isIPOOpen checks if an IPO is currently open
func (f *IPOFetcher) isIPOOpen(ipo models.IPO, now time.Time) bool {
	// Parse dates (assuming format like "15 Jan 2024" or "15-01-2024")
	openDate, err := f.parseDate(ipo.OpenDate)
	if err != nil {
		return false
	}

	closeDate, err := f.parseDate(ipo.CloseDate)
	if err != nil {
		return false
	}

	// Check if current date is between open and close dates
	return (now.After(openDate) || now.Equal(openDate)) && (now.Before(closeDate) || now.Equal(closeDate))
}

// parseDate parses various date formats commonly used in Indian IPO listings
func (f *IPOFetcher) parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	// Clean up common prefixes/suffixes
	dateStr = strings.TrimPrefix(dateStr, "Opens:")
	dateStr = strings.TrimPrefix(dateStr, "Closes:")
	dateStr = strings.TrimSpace(dateStr)

	// Try different date formats commonly used in Indian websites
	formats := []string{
		"02 Jan 2006",
		"02-Jan-2006",
		"02/01/2006",
		"02-01-2006",
		"2006-01-02",
		"Jan 02, 2006",
		"January 02, 2006",
		"02 January 2006",
		"02 Jan 06",
		"02-Jan-06",
		"02/01/06",
	}

	// Try parsing in IST timezone first
	ist, _ := time.LoadLocation("Asia/Kolkata")
	for _, format := range formats {
		if t, err := time.ParseInLocation(format, dateStr, ist); err == nil {
			return t, nil
		}
		// Also try without timezone
		if t, err := time.Parse(format, dateStr); err == nil {
			// Convert to IST
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ist), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

