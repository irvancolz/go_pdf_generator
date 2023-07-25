package db

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("failed to create mock :", err)
	}

	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "John Doe", "john.doe@example.com")
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").WithArgs(1).WillReturnRows(rows)
	var user User
	err = sqlxDB.Get(&user, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting user by id", err)
	}
	if user.Name != "John Doe" {
		t.Errorf("expected name to be %s but got %s", "John Doe", user.Name)
	}
	if user.Email != "john.doe@example.com" {
		t.Errorf("expected email to be %s but got %s", "john.doe@example.com", user.Email)
	}

}

func ConnectPostgres() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", "host=localhost port=5432 user=postgres password=root dbname=test sslmode=disable")
	if err != nil {
		log.Println("failed to connect to database :", err)
		return nil, err
	}
	defer db.Close()
	return db, nil
}

// func ConnectMock() (*sqlx.DB, error) {

// }

func TestConnection(t *testing.T) {

}
