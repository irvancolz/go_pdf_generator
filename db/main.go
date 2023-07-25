package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	// _ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "host=localhost port=5432 user=postgres password=root dbname=test sslmode=disable")
	if err != nil {
		log.Println("failed to connect to database :", err)
		return
	}
	defer db.Close()

}
