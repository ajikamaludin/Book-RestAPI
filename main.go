package main

import (
	"grace/database"
	"grace/handlers"
	"grace/repositories"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Connect to the database
	db := setupDatabase()
	migrateDatabase(db)

	// Initialize repositories
	bookRepository := &repositories.BookRepository{DB: db}
	collectionRepository := &repositories.CollectionRepository{DB: db}
	collectionBookRepository := &repositories.CollectionBookRepository{DB: db}

	// Initialize handlers
	bookHandler := &handlers.BookHandler{Repository: bookRepository}
	collectionHandler := &handlers.CollectionHandler{
		Repository:         collectionRepository,
		BookRepository:     bookRepository,
		CollectionBookRepo: collectionBookRepository,
	}

	// Create Fiber app
	app := fiber.New()

	// Book endpoints
	app.Post("/books", bookHandler.CreateBook)
	app.Put("/books/:id", bookHandler.UpdateBook)
	app.Delete("/books/:id", bookHandler.DeleteBook)
	app.Get("/books", bookHandler.GetBooks)
	app.Get("/books/:id", bookHandler.GetBook)

	// Collection endpoints
	app.Post("/collections", collectionHandler.CreateCollection)
	app.Put("/collections/:id", collectionHandler.UpdateCollection)
	app.Delete("/collections/:id", collectionHandler.DeleteCollection)
	app.Get("/collections/:id", collectionHandler.GetCollection)
	app.Get("/collections/:id/books", collectionHandler.GetCollectionBooks)                  // New route for getting collection books
	app.Post("/collections/:id/books", collectionHandler.AddBookToCollection)                // New route for adding book to collection
	app.Delete("/collections/:id/books/:bookId", collectionHandler.RemoveBookFromCollection) // New route for removing book from collection

	// Start the server
	app.Listen(":3000")
}

func setupDatabase() *gorm.DB {
	dsn := "host=localhost user=app password=password dbname=app port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func migrateDatabase(db *gorm.DB) {
	db.AutoMigrate(&database.Book{}, &database.Collection{}, &database.CollectionBook{})
}
