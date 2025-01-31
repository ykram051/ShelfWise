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
	bookRepo := repositories.NewInMemoryBookStore("./books.json")
	authorRepo := repositories.NewInMemoryAuthorStore("./authors.json")
	customerRepo := repositories.NewInMemoryCustomerStore("./customers.json")
	orderRepo := repositories.NewInMemoryOrderStore("./orders.json")

	// Services
	bookService := services.NewBookService(bookRepo, authorRepo)
	reportService := services.NewReportService(orderRepo)
	authorService := services.NewAuthorService(authorRepo, bookRepo)
	customerService := services.NewCustomerService(customerRepo)
	orderService := services.NewOrderService(orderRepo, bookRepo, customerRepo)

	// Controllers
	bookController := controllers.NewBookController(bookService)
	authorController := controllers.NewAuthorController(authorService)
	customerController := controllers.NewCustomerController(customerService)
	orderController := controllers.NewOrderController(orderService)
	reportController := controllers.NewReportController(reportService)

	views.StartDailyReportJob(reportService)

	// Books
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
	// Authors
	http.HandleFunc("/authors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			authorController.CreateAuthor(w, r)
		case http.MethodGet:
			authorController.GetAuthor(w, r)
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

	// Customers
	http.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			customerController.CreateCustomer(w, r)
		case http.MethodGet:
			customerController.GetCustomer(w, r)
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

	// Orders
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderController.CreateOrder(w, r)
		case http.MethodGet:
			orderController.GetOrder(w, r)
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
	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			reportController.ListReports(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on 8086")
	if err := http.ListenAndServe(":8086", nil); err != nil {
		log.Fatal(err)
	}

	bookRepo.Save()
	authorRepo.Save()
	customerRepo.Save()
	orderRepo.Save()
}
