# Bookstore API - README

## Overview
The Bookstore API is a RESTful service built using Go, designed to manage books, authors, customers, orders, and sales reports. This README provides instructions on setting up and using the API.

## Prerequisites
To run the project, ensure you have the following installed:

- **Go**.

## Setup Instructions

### 1. Clone the Repository
Clone the project repository to your local machine:

```bash
git clone https://github.com/ykram051/Intro-to-Go-.git
cd '.\Final Project\'
```

### 3. Install Dependencies
Use `go mod` to install project dependencies:


### 4. Run the Server
Run the API server with:

```bash
go run main.go
```

The server will start on the port `8086`.

### 5. Test the API
Use the following base URL for requests:

```
http://localhost:8086
```

You can use Postman to test endpoints. For example, to list all books:

```bash
GET http://localhost:8086/books
```

## Key Endpoints
Here are the main API endpoints:

### Authors
- **GET /authors**: List all authors or fetch an author by ID.
- **POST /authors**: Create a new author.
- **PUT /authors/{id}**: Update author details.
- **DELETE /authors/{id}**: Delete an author (if no books are associated).

### Books
- **GET /books**: List all books or search by criteria (title, author, genre).
- **POST /books**: Add a new book.
- **PUT /books/{id}**: Update book details.
- **DELETE /books/{id}**: Delete a book.

### Customers
- **GET /customers**: List all customers or fetch by ID.
- **POST /customers**: Add a new customer.
- **PUT /customers/{id}**: Update customer details.
- **DELETE /customers/{id}**: Delete a customer.

### Orders
- **GET /orders**: List all orders or filter by date range.
- **POST /orders**: Create a new order.
- **DELETE /orders/{id}**: Cancel an order.

### Reports
- **GET /report**: Retrieve sales reports for a specified date range.

## Development Notes

### Project Structure
The project follows a modular structure:

- **models/**: Data models for books, authors, customers, orders, etc.
- **controllers/**: Handle HTTP requests and responses.
- **services/**: Contain business logic.
- **repositories/**: Handle data persistence (in-memory or database).
- **main.go**: Entry point for the application.


### Error Handling
Errors are returned in a log file named error.log

## Feedback and Issues
For questions or issues, you can contact me via email in : ikram.benfellah@um6p.ma .
