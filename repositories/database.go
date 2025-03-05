package repositories

import (
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func InitDB() {
	dsn := "postgres://postgres:root@localhost:5432/shelfwise?sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	if err := sqldb.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	DB = bun.NewDB(sqldb, pgdialect.New())
	log.Println("Connected to PostgreSQL!")
}

// Close DB Connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
