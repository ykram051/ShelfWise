package main

import (
	"FinalProject/controllers"
	"FinalProject/middleware"
	"FinalProject/repositories"
	"FinalProject/services"
	"FinalProject/task"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // Import the godotenv package
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	repositories.InitDB()
	defer repositories.CloseDB()

	// Initialize repositories
	authorRepo := repositories.NewAuthorRepository(repositories.DB)
	bookRepo := repositories.NewBookRepository(repositories.DB)
	customerRepo := repositories.NewCustomerRepository(repositories.DB)
	orderRepo := repositories.NewOrderRepository(repositories.DB)
	reportRepo := repositories.NewReportStore(repositories.DB)
	userRepo := repositories.NewUserRepository(repositories.DB)

	// Initialize services
	authorService := services.NewAuthorService(authorRepo)
	bookService := services.NewBookService(bookRepo, authorRepo)
	customerService := services.NewCustomerService(customerRepo)
	orderService := services.NewOrderService(orderRepo, bookRepo, customerRepo)
	reportService := services.NewReportService(orderRepo, reportRepo)
	authService := services.NewAuthService(userRepo)

	// Initialize controllers
	authorController := controllers.NewAuthorController(authorService)
	bookController := controllers.NewBookController(bookService)
	customerController := controllers.NewCustomerController(customerService)
	orderController := controllers.NewOrderController(orderService)
	reportController := controllers.NewReportController(reportService)
	authController := controllers.NewAuthController(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, orderService)

	// Start background tasks
	task.StartDailyReportJob(reportService)

	// Setup router
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/register", authController.Register).Methods("POST")
	router.HandleFunc("/login", authController.Login).Methods("POST")

	// Protected API routes (JWT required)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware.JWTAuthMiddleware)

	// 📚 Book routes
	api.HandleFunc("/books", bookController.CreateBook).Methods("POST")
	api.HandleFunc("/books", bookController.SearchBooks).Methods("GET")
	api.HandleFunc("/books/{id:[0-9]+}", bookController.GetBook).Methods("GET")
	api.HandleFunc("/books/{id}", bookController.UpdateBook).Methods("PUT")
	api.HandleFunc("/books/{id}", bookController.DeleteBook).Methods("DELETE")

	// ✍️ Author routes
	api.HandleFunc("/authors", authorController.CreateAuthor).Methods("POST")
	api.HandleFunc("/authors", authorController.SearchAuthors).Methods("GET")
	api.HandleFunc("/authors/{id}", authorController.GetAuthor).Methods("GET")
	api.HandleFunc("/authors/{id}", authorController.UpdateAuthor).Methods("PUT")
	api.HandleFunc("/authors/{id}", authorController.DeleteAuthor).Methods("DELETE")

	// 👥 Customer routes
	api.HandleFunc("/customers", customerController.ListCustomers).Methods("GET")
	api.HandleFunc("/customers/{id}", customerController.GetCustomer).Methods("GET")
	api.HandleFunc("/customers/{id}", customerController.UpdateCustomer).Methods("PUT")
	api.HandleFunc("/customers/{id}", customerController.DeleteCustomer).Methods("DELETE")

	// 📦 Order routes
	api.HandleFunc("/orders", orderController.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", orderController.ListOrders).Methods("GET")
	api.HandleFunc("/orders/{id}", orderController.GetOrder).Methods("GET")
	api.HandleFunc("/orders/{id}", orderController.UpdateOrder).Methods("PUT")
	api.HandleFunc("/orders/{id}", orderController.DeleteOrder).Methods("DELETE")
	api.HandleFunc("/orders/date-range", orderController.GetOrdersByDateRange).Methods("GET")
	api.HandleFunc("/orders/search-by-customer", orderController.SearchOrdersByCustomerID).Methods("GET")

	// 📊 Report routes
	api.HandleFunc("/report", reportController.ListReports).Methods("GET")

	// Start server
	port := os.Getenv("PORT") // Get the port from the environment variables
	if port == "" {
		port = "8086" // Default port if not specified
	}
	log.Printf("Server running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server error:", err)
	}
}
