package controllers

import (
	"FinalProject/models"
	"FinalProject/services"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController(rs *services.ReportService) *ReportController {
	return &ReportController{reportService: rs}
}

func (rc *ReportController) ListReports(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error
	if fromStr == "" {
		from = time.Date(2025, 1, 1, 15, 0, 0, 0, time.UTC) // Changed to 15:00
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, "invalid 'from' date format (use YYYY-MM-DD)")
			return
		}
		from = from.UTC()
	}

	if toStr == "" {
		to = time.Now().Truncate(24 * time.Hour).Add(24*time.Hour - time.Second).UTC()
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, "invalid 'to' date format (use YYYY-MM-DD)")
			return
		}
		to = to.UTC()
	}

	reportsDir, err := filepath.Abs(filepath.Join("..", "Final Project", "output-reports"))
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("failed to resolve reports directory: %v", err))
		return
	}
	fmt.Println("Resolved reports directory:", reportsDir)

	files, err := os.ReadDir(reportsDir)
	if os.IsNotExist(err) {
		fmt.Println("Reports directory does not exist:", reportsDir)
		_ = json.NewEncoder(w).Encode([]models.SalesReport{})
		return
	} else if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("failed to read reports directory: %v", err))
		return
	}

	var matchingReports []models.SalesReport

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if !strings.HasPrefix(filename, "report_") || !strings.HasSuffix(filename, ".json") {
			fmt.Println("Skipping non-report file:", filename)
			continue
		}

		trimmed := strings.TrimPrefix(filename, "report_")
		trimmed = strings.TrimSuffix(trimmed, ".json")
		fileTimestamp, parseErr := time.Parse("010220061504", trimmed) // Adjusted timestamp format
		if parseErr != nil {
			fmt.Println("Failed to parse timestamp for file:", filename, "Error:", parseErr)
			continue
		}

		fileTimestamp = fileTimestamp.UTC()

		if (fileTimestamp.After(from) || fileTimestamp.Equal(from)) && (fileTimestamp.Before(to) || fileTimestamp.Equal(to)) {
			path := filepath.Join(reportsDir, filename)
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				fmt.Println("Failed to read file:", path, "Error:", readErr)
				continue
			}

			var sr models.SalesReport
			if unmarshalErr := json.Unmarshal(data, &sr); unmarshalErr != nil {
				fmt.Println("Failed to unmarshal file:", path, "Error:", unmarshalErr)
				continue
			}

			matchingReports = append(matchingReports, sr)
		} else {
			fmt.Printf("File outside date range: %s (Timestamp: %v, From: %v, To: %v)\n", filename, fileTimestamp, from, to)
		}

	}

	fmt.Printf("Number of matching reports: %d\n", len(matchingReports))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matchingReports)
}
