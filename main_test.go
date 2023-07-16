package main_test

import (
	"bytes"
	"encoding/json"
	"grace/handlers"
	"grace/repositories"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCreateBookEndpoint(t *testing.T) {
	// Create a sample book payload
	book := map[string]interface{}{
		"id":           9999,
		"title":        "Test Book",
		"author":       "John Doe",
		"published_at": "2023-01-01T00:00:00Z",
		"edition":      "First Edition",
		"description":  "This is a test book",
		"genre":        "Fiction",
	}
	payload, _ := json.Marshal(book)

	// Create a new Fiber app
	app := fiber.New()

	// Register the route handler for creating a book
	db := setupDatabase()
	bookRepository := &repositories.BookRepository{DB: db}
	bookHandler := &handlers.BookHandler{Repository: bookRepository}
	app.Post("/books", bookHandler.CreateBook)

	// Create a request to the create book endpoint
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	resp, _ := app.Test(req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestListBooksEndpoint(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Register the route handler for listing books
	db := setupDatabase()
	bookRepository := &repositories.BookRepository{DB: db}
	bookHandler := &handlers.BookHandler{Repository: bookRepository}
	app.Get("/books", bookHandler.GetBooks)

	// Create a request to the list books endpoint
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	resp, _ := app.Test(req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEditBookEndpoint(t *testing.T) {
	// Create a sample book payload for editing
	book := map[string]interface{}{
		"title": "Updated Book Title",
	}

	// Marshal the book payload to JSON
	payload, _ := json.Marshal(book)

	// Create a new Fiber app
	app := fiber.New()

	// Register the route handler for editing a book
	db := setupDatabase()
	bookRepository := &repositories.BookRepository{DB: db}
	bookHandler := &handlers.BookHandler{Repository: bookRepository}
	app.Put("/books/:id", bookHandler.UpdateBook)

	// Create a request to the edit book endpoint
	req := httptest.NewRequest(http.MethodPut, "/books/9999", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	resp, _ := app.Test(req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// TODO: Validate the response body or check the database for the updated book
}

func TestDeleteBookEndpoint(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Register the route handler for deleting a book
	db := setupDatabase()
	bookRepository := &repositories.BookRepository{DB: db}
	bookHandler := &handlers.BookHandler{Repository: bookRepository}
	app.Delete("/books/:id", bookHandler.DeleteBook)

	// Create a request to the delete book endpoint
	req := httptest.NewRequest(http.MethodDelete, "/books/9999", nil)
	resp, _ := app.Test(req)

	// Check the response status code
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func setupDatabase() *gorm.DB {
	dsn := "host=localhost user=app password=password dbname=app port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}
