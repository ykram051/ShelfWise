package repositories

import (
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// DB instance
var DB *bun.DB

// Initialize PostgreSQL connection
func InitDB() {
	dsn := "postgres://postgres:root@localhost:5432/shelfwise?sslmode=disable"

	// Open PostgreSQL database connection
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := sqldb.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// ✅ Use `pgdialect.New()` for PostgreSQL
	DB = bun.NewDB(sqldb, pgdialect.New())
	log.Println("Connected to PostgreSQL!")
}

// Close DB Connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
