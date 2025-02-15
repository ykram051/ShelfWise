package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"context"
	"time"
)

type ReportService struct {
	orderStore  repositories.OrderStore
	reportStore repositories.ReportStore // ✅ Add a store for reports
}

func NewReportService(os repositories.OrderStore, rs repositories.ReportStore) *ReportService {
	return &ReportService{orderStore: os, reportStore: rs}
}

func (rs *ReportService) GenerateSalesReport(ctx context.Context, from, to time.Time) (models.SalesReport, error) {
	select {
	case <-ctx.Done():
		return models.SalesReport{}, ctx.Err()
	default:
	}

	orders, err := rs.orderStore.GetOrdersByDateRange(from, to)
	if err != nil {
		return models.SalesReport{}, err
	}

	totalRevenue := 0.0
	totalOrders := len(orders)
	bookSalesMap := make(map[int]int)
	bookMap := make(map[int]*models.Book) // ✅ Store pointers instead of values

	for _, o := range orders {
		totalRevenue += o.TotalPrice
		for _, item := range o.Items {
			bookSalesMap[item.BookID] += item.Quantity
			bookMap[item.BookID] = item.Book // ✅ Store the pointer directly
		}
	}
	var allbooks []models.BookSales
	for bookID, qty := range bookSalesMap {
		
		allbooks = append(allbooks, models.BookSales{
			BookID:   bookID,
			Book:     bookMap[bookID], // ✅ Now bookMap[bookID] is a pointer (*models.Book)
			Quantity: qty,
		})
	}


	var topSelling []models.BookSales
	for bookID, qty := range bookSalesMap {
		
		topSelling = append(topSelling, models.BookSales{
			BookID:   bookID,
			Book:     bookMap[bookID], // ✅ Now bookMap[bookID] is a pointer (*models.Book)
			Quantity: qty,
		})
	}

	report := models.SalesReport{
		Timestamp:       time.Now(),
		TotalRevenue:    totalRevenue,
		TotalOrders:     totalOrders,
		TopSellingBooks: topSelling,
	}

	// ✅ Insert the report into the database
	err = rs.reportStore.SaveReport(ctx, &report)
	if err != nil {
		return models.SalesReport{}, err
	}

	return report, nil
}
