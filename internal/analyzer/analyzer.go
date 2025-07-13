package analyzer

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type AnalysisResult struct {
	HTMLVersion   string
	Title         string
	Headings      map[string]int
	InternalLinks []string
	ExternalLinks []string
	BrokenLinks   []string
	HasLoginForm  bool
}

func AnalyzeURL(target string) (AnalysisResult, error) {
	resp, err := http.Get(target)
	if err != nil {
		return AnalysisResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return AnalysisResult{}, errors.New("failed to fetch target URL")
	}

	z := html.NewTokenizer(resp.Body)
	result := AnalysisResult{
		Headings:      map[string]int{},
		InternalLinks: []string{},
		ExternalLinks: []string{},
		BrokenLinks:   []string{},
	}

	parsedBase, err := url.Parse(target)
	if err != nil {
		return result, err
	}

	var inTitle bool
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// End of document
			return result, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()

			switch t.Data {
			case "html":
				for _, attr := range t.Attr {
					if attr.Key == "version" {
						result.HTMLVersion = attr.Val
					}
				}
				if result.HTMLVersion == "" {
					result.HTMLVersion = "HTML5"
				}
			case "title":
				inTitle = true
			case "h1", "h2", "h3", "h4", "h5", "h6":
				result.Headings[strings.ToLower(t.Data)]++
			case "a":
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						link := attr.Val
						if strings.HasPrefix(link, "http") {
							if strings.Contains(link, parsedBase.Host) {
								result.InternalLinks = append(result.InternalLinks, link)
							} else {
								result.ExternalLinks = append(result.ExternalLinks, link)
							}
						} else if strings.HasPrefix(link, "/") {
							result.InternalLinks = append(result.InternalLinks, parsedBase.Scheme+"://"+parsedBase.Host+link)
						}
					}
				}
			case "input":
				for _, attr := range t.Attr {
					if attr.Key == "type" && attr.Val == "password" {
						result.HasLoginForm = true
					}
				}
			}
		case html.TextToken:
			if inTitle {
				result.Title = strings.TrimSpace(z.Token().Data)
				inTitle = false
			}
		}
	}
}
