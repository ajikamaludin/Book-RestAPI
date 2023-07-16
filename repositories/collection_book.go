package repositories

import (
	"grace/database"

	"gorm.io/gorm"
)

type CollectionBookRepository struct {
	DB *gorm.DB
}

func (r *CollectionBookRepository) AddBookToCollection(book *database.Book, collection *database.Collection) error {
	return r.DB.Model(collection).Association("Books").Append(book)
}

func (r *CollectionBookRepository) RemoveBookFromCollection(book *database.Book, collection *database.Collection) error {
	return r.DB.Model(collection).Association("Books").Delete(book)
}

func (r *CollectionBookRepository) RemoveCollectionBook(book *database.Book, collection *database.Collection) error {
	return r.DB.Delete(&database.CollectionBook{CollectionID: collection.ID, BookID: book.ID}).Error
}
