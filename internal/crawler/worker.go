package crawler

import (
	"log"
	"time"

	"github.com/adithimanjunath/crawler-api/internal/analyzer"
	"github.com/adithimanjunath/crawler-api/internal/db"
	"github.com/adithimanjunath/crawler-api/internal/models"
	"gorm.io/gorm/clause"
)

// StartCrawlWorker launches the background goroutine to process queued URLs
func StartCrawlWorker() {
	go func() {
		for {
			processQueuedURLs()
			time.Sleep(10 * time.Second)
		}
	}()
}

func processQueuedURLs() {
	var urls []models.URL

	// Find all queued URLs
	if err := db.DB.Where("status = ?", "queued").Find(&urls).Error; err != nil {
		log.Printf("Error fetching queued URLs: %v", err)
		return
	}

	for _, url := range urls {
		log.Printf("Analyzing URL: %s", url.URL)

		// Mark as running
		url.Status = "running"
		db.DB.Save(&url)

		// Analyze the URL
		result, err := analyzer.AnalyzeURL(url.URL)
		if err != nil {
			log.Printf("Error analyzing %s: %v", url.URL, err)
			url.Status = "error"
			db.DB.Save(&url)
			continue
		}

		// Prepare analysis result
		newData := models.AnalysisResult{
			URLID:              url.ID,
			HTMLVersion:        result.HTMLVersion,
			Title:              result.Title,
			H1Count:            result.Headings["h1"],
			H2Count:            result.Headings["h2"],
			H3Count:            result.Headings["h3"],
			H4Count:            result.Headings["h4"],
			H5Count:            result.Headings["h5"],
			H6Count:            result.Headings["h6"],
			InternalLinksCount: len(result.InternalLinks),
			ExternalLinksCount: len(result.ExternalLinks),
			BrokenLinksCount:   len(result.BrokenLinks),
			HasLoginForm:       result.HasLoginForm,
		}

		// UPSERT using GORM clause.OnConflict
		err = db.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url_id"}}, // unique key column
			UpdateAll: true,                              // update all fields if conflict
		}).Create(&newData).Error

		if err != nil {
			log.Printf("Failed to upsert analysis for %s: %v", url.URL, err)
		}

		// Mark as done
		url.Status = "done"
		db.DB.Save(&url)

		log.Printf("Analysis complete for: %s", url.URL)
	}
}
