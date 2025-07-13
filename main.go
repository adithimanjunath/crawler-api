// main.go
package main

import (
	"log"
	"os"
	"time"

	"github.com/adithimanjunath/crawler-api/internal/crawler"
	"github.com/adithimanjunath/crawler-api/internal/db"
	"github.com/adithimanjunath/crawler-api/internal/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	db.Connect()
	db.DB.AutoMigrate(&models.URL{}, &models.AnalysisResult{})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/urls", func(c *gin.Context) {
		var url models.URL
		if err := c.ShouldBindJSON(&url); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		url.Status = "queued"
		if err := db.DB.Create(&url).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to save URL"})
			return
		}
		c.JSON(201, url)
	})

	r.GET("/urls", func(c *gin.Context) {
		var urls []models.URL
		if err := db.DB.Find(&urls).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch URLs"})
			return
		}
		c.JSON(200, urls)
	})

	r.GET("/urls/:id/results", func(c *gin.Context) {
		id := c.Param("id")
		var result models.AnalysisResult
		if err := db.DB.Where("url_id = ?", id).First(&result).Error; err != nil {
			c.JSON(404, gin.H{"error": "No analysis result found"})
			return
		}

		headings := map[string]int{
			"h1": result.H1Count,
			"h2": result.H2Count,
			"h3": result.H3Count,
			"h4": result.H4Count,
			"h5": result.H5Count,
			"h6": result.H6Count,
		}

		c.JSON(200, gin.H{
			"url_id":         result.URLID,
			"html_version":   result.HTMLVersion,
			"title":          result.Title,
			"headings":       headings,
			"internal_links": result.InternalLinksCount,
			"external_links": result.ExternalLinksCount,
			"broken_links":   result.BrokenLinksCount,
			"has_login_form": result.HasLoginForm,
		})
	})

	r.DELETE("/urls/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.DB.Where("url_id = ?", id).Delete(&models.AnalysisResult{})
		if err := db.DB.Delete(&models.URL{}, id).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete URL"})
			return
		}
		c.Status(204)
	})

	r.POST("/urls/:id/reanalyze", func(c *gin.Context) {
		id := c.Param("id")
		var url models.URL
		if err := db.DB.First(&url, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "URL not found"})
			return
		}
		url.Status = "queued"
		if err := db.DB.Save(&url).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to re-queue URL"})
			return
		}
		c.JSON(200, gin.H{"message": "Re-analysis started"})
	})

	crawler.StartCrawlWorker()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
