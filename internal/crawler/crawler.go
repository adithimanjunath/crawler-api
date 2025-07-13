package crawler

import (
	"net/http"
)

// CrawlResult holds extracted data from a URL
type CrawlResult struct {
	HTMLVersion   string
	Title         string
	Headings      map[string]int
	InternalLinks []string
	ExternalLinks []string
	BrokenLinks   []BrokenLink
	HasLoginForm  bool
}

// BrokenLink holds info about a failed request
type BrokenLink struct {
	URL    string
	Status int
}

// AnalyzeURL crawls and analyzes the given URL

func checkLink(url string) int {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return 500
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
