package repositories

import (
	"FinalProject/models"
	"context"

	"github.com/uptrace/bun"
)

type ReportStore struct {
	db *bun.DB
}

func NewReportStore(db *bun.DB) ReportStore {
	return ReportStore{db: db}
}

// ✅ Save report to the database
func (rs *ReportStore) SaveReport(ctx context.Context, report *models.SalesReport) error {
	_, err := rs.db.NewInsert().
		Model(report).
		Returning("id"). // ✅ Get the generated ID
		Exec(ctx)

	return err
}
