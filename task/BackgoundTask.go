package task

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

// StartDailyReportJob schedules report generation every day at 00:00 UTC
func StartDailyReportJob(rs *services.ReportService) {
	go func() {
		for {
			now := time.Now().UTC().Truncate(time.Second) // Ensure precision

			// ✅ Schedule next run for 00:00 (midnight) UTC
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

			// ✅ If it's already past 00:00 today, schedule for tomorrow
			if now.After(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			timeUntilNextRun := time.Until(nextRun)
			log.Printf("Next report scheduled at: %v (in %v)", nextRun, timeUntilNextRun)

			time.Sleep(timeUntilNextRun)

			// ✅ Adjusting timestamps to cover the full previous day (00:00 - 23:59:59)
			from := nextRun.Add(-24 * time.Hour).UTC().Truncate(time.Second)
			to := nextRun.UTC().Truncate(time.Second)

			log.Printf("Generating report for range: %v to %v", from, to)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			report, err := rs.GenerateSalesReport(ctx, from, to)
			if err != nil {
				controllers.LogError(err)
				continue
			}

			// ✅ Ensure orders are detected correctly
			if report.TotalOrders == 0 {
				log.Println("⚠ No new orders found in the time range. Skipping report generation.")
				continue
			}

			// Ensure reports directory exists
			outputDir := "output-reports"
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				controllers.LogError(err)
				continue
			}

			// Generate filename with timestamp
			filename := "report_" + report.Timestamp.Format("010220060000") + ".json" // Fixed at 00:00
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

			log.Printf("✅ Sales report generated: %s with ID: %d", path, report.ID)
		}
	}()
}
