package services

import (
	"FinalProject/models"
	"FinalProject/repositories"
	"time"
)

type ReportService struct {
	orderStore repositories.OrderStore
}

func NewReportService(os repositories.OrderStore) *ReportService {
	return &ReportService{orderStore: os}
}

func (rs *ReportService) GenerateSalesReport(from, to time.Time) (models.SalesReport, error) {
	orders, err := rs.orderStore.GetOrdersByDateRange(from, to)
	if err != nil {
		return models.SalesReport{}, err
	}

	totalRevenue := 0.0
	totalOrders := len(orders)
	bookSalesMap := make(map[int]int)
	bookMap := make(map[int]models.Book)

	for _, o := range orders {
		totalRevenue += o.TotalPrice
		for _, item := range o.Items {
			bookSalesMap[item.Book.ID] += item.Quantity
			bookMap[item.Book.ID] = item.Book
		}
	}

	var topSelling []models.BookSales
	for bookID, qty := range bookSalesMap {
		topSelling = append(topSelling, models.BookSales{
			Book:     bookMap[bookID],
			Quantity: qty,
		})
	}

	report := models.SalesReport{
		Timestamp:       time.Now(),
		TotalRevenue:    totalRevenue,
		TotalOrders:     totalOrders,
		TopSellingBooks: topSelling,
	}
	return report, nil
}
