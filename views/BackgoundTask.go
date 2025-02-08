package views

import (
	"FinalProject/controllers"
	"FinalProject/services"
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

// StartDailyReportJob schedules a daily report generation job
func StartDailyReportJob(rs *services.ReportService) {
	go func() {
		for {
			now := time.Now().UTC()
			nextRun := now.Truncate(24 * time.Hour).Add(24 * time.Hour)

			timeUntilNextRun := time.Until(nextRun)
			log.Printf("Next daily report scheduled in: %v", timeUntilNextRun)

			time.Sleep(timeUntilNextRun)

			from := now.Truncate(24 * time.Hour).Add(-24 * time.Hour)
			to := now.Truncate(24 * time.Hour)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			report, err := rs.GenerateSalesReport(ctx, from, to)
			if err != nil {
				controllers.LogError(err)
				continue
			}

			// Ensure reports directory exists
			outputDir := "output-reports"
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				controllers.LogError(err)
				continue
			}

			// Generate filename with timestamp
			filename := "report_" + report.Timestamp.Format("20060102_150405") + ".json"
			path := filepath.Join(outputDir, filename)

			// Marshal report data
			data, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				controllers.LogError(err)
				continue
			}

			// Write report to file
			if err := os.WriteFile(path, data, 0644); err != nil {
				controllers.LogError(err)
				continue
			}

			log.Println("✅ Daily sales report generated:", path)
		}
	}()
}
