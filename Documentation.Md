# Project Documentation: Comprehensive Overview

## Introduction

This project is a RESTful web service designed to manage books, authors, customers, orders, and reports. The application provides endpoints for creating, reading, updating, and deleting resources while maintaining data integrity through robust validation and service-layer logic.

## Key Components and Logic

### 1. Author Management

#### Endpoints:
- **POST /authors**: Create a new author.
- **GET /authors**: List all authors or fetch details by ID.
- **PUT /authors/{id}**: Update author details.
- **DELETE /authors/{id}**: Delete an author.

#### Logic:
- Authors are stored in an in-memory store.
- Deleting an author is only allowed if no books are associated with them.
- Before deleting, a check is performed using `bookStore.ListBooks()` to find any associated books.

### 2. Book Management

#### Endpoints:
- **POST /books**: Add a new book.
- **GET /books**: List all books or fetch details by ID.
- **PUT /books/{id}**: Update book details.
- **DELETE /books/{id}**: Delete a book.

#### Logic:
- Books are linked to authors by their `author.id`.
- Before adding a book, the existence of the specified author is validated, if the author doesn't exist it automatically creates it.
- Stock is adjusted dynamically during order creation.

### 3. Customer Management

#### Endpoints:
- **POST /customers**: Add a new customer.
- **GET /customers**: List all customers or fetch details by ID.
- **PUT /customers/{id}**: Update customer details.
- **DELETE /customers/{id}**: Delete a customer.

#### Logic:
- Customers are stored with essential details such as name, email, and address.


### 4. Order Management

#### Endpoints:
- **POST /orders**: Create a new order.
- **GET /orders**: List all orders or fetch details by ID.
- **PUT /orders/{id}**: Update an order with the given id.
- **DELETE /orders/{id}**: Cancel an order.

#### Logic:
- Orders include customer details and the books being purchased.
- Stock is verified before creating an order to ensure availability.
- Total price is calculated based on the quantity and price of each book.
- Orders cannot be created for nonexistent customers , or nonexistent books.
- When updating or deleting and order ,we make sure the book stocks are updated as well.

### 5. Report Management

#### Endpoints:
- **GET /reports**: Retrieve sales reports within a specified date range.

#### Logic:
- A background job generates daily sales reports at midnight, summarizing:
  - Total revenue.
  - Total number of orders.
  - Total books sold.
  - Top-selling books.
- Reports are stored in JSON files within the `output-reports` directory.

## Service and Store Layer Design

### Service Layer:
- Contains business logic and validations.
- Coordinates actions between stores to maintain data integrity.

### Store Layer:
- Handles CRUD operations for in-memory data storage.
- Ensures thread safety using mutex locks.
