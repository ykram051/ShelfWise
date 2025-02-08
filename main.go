package main

import (
	"FinalProject/controllers"
	"FinalProject/repositories"
	"FinalProject/services"
	"FinalProject/views"
	"log"
	"net/http"
)

func main() {
	// Initialize PostgreSQL Database
	repositories.InitDB()
	defer repositories.CloseDB()

	// Initialize repositories with Bun DB
	authorRepo := repositories.NewAuthorRepository(repositories.DB)
	bookRepo := repositories.NewBookRepository(repositories.DB)
	customerRepo := repositories.NewCustomerRepository(repositories.DB)
	orderRepo := repositories.NewOrderRepository(repositories.DB)

	// Initialize services
	authorService := services.NewAuthorService(authorRepo)
	bookService := services.NewBookService(bookRepo, authorRepo)
	customerService := services.NewCustomerService(customerRepo)
	orderService := services.NewOrderService(orderRepo, bookRepo, customerRepo)
	reportService := services.NewReportService(orderRepo)

	// Initialize controllers
	authorController := controllers.NewAuthorController(authorService)
	bookController := controllers.NewBookController(bookService)
	customerController := controllers.NewCustomerController(customerService)
	orderController := controllers.NewOrderController(orderService)
	reportController := controllers.NewReportController(reportService)

	// Start Daily Report Job
	views.StartDailyReportJob(reportService)

	// Books Routes
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			bookController.CreateBook(w, r)
		case http.MethodGet:
			bookController.SearchBooks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookController.GetBook(w, r)
		case http.MethodPut:
			bookController.UpdateBook(w, r)
		case http.MethodDelete:
			bookController.DeleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Authors Routes
	http.HandleFunc("/authors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			authorController.CreateAuthor(w, r)
		case http.MethodGet:
			authorController.ListAuthors(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/authors/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authorController.GetAuthor(w, r)
		case http.MethodPut:
			authorController.UpdateAuthor(w, r)
		case http.MethodDelete:
			authorController.DeleteAuthor(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Customers Routes
	http.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			customerController.CreateCustomer(w, r)
		case http.MethodGet:
			customerController.ListCustomers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			customerController.GetCustomer(w, r)
		case http.MethodPut:
			customerController.UpdateCustomer(w, r)
		case http.MethodDelete:
			customerController.DeleteCustomer(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Orders Routes
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderController.CreateOrder(w, r)
		case http.MethodGet:
			orderController.ListOrders(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderController.GetOrder(w, r)
		case http.MethodPut:
			orderController.UpdateOrder(w, r)
		case http.MethodDelete:
			orderController.DeleteOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Reports Route
	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			reportController.ListReports(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("✅ Server running on port 8086...")
	if err := http.ListenAndServe(":8086", nil); err != nil {
		log.Fatal(err)
	}
}
