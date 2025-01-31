package views

import (
	"FinalProject/controllers"
	"FinalProject/services"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

func StartDailyReportJob(rs *services.ReportService) {
	go func() {
		for {
			now := time.Now()
			nextRun := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
			if now.After(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			timeUntilNextRun := nextRun.Sub(now)
			time.Sleep(timeUntilNextRun)

			from := now.Truncate(24 * time.Hour).Add(-24 * time.Hour)
			to := now.Truncate(24 * time.Hour)

			report, err := rs.GenerateSalesReport(from, to)
			if err != nil {
				controllers.LogError(err)
				continue
			}

			os.MkdirAll("output-reports", os.ModePerm)
			filename := "report_" + report.Timestamp.Format("01022006150405") + ".json"
			path := filepath.Join("output-reports", filename)

			data, err := json.MarshalIndent(report, "", "  ")
			if err != nil {
				controllers.LogError(err)
				continue
			}

			if err := os.WriteFile(path, data, 0644); err != nil {
				controllers.LogError(err)
				continue
			}

			log.Println("Daily report generated:", path)
		}
	}()
}
