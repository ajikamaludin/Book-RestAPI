package repositories

import (
	"grace/database"
	"time"

	"gorm.io/gorm"
)

type BookRepository struct {
	DB *gorm.DB
}

func (r *BookRepository) Create(book *database.Book) error {
	return r.DB.Create(book).Error
}

func (r *BookRepository) Update(book *database.Book) error {
	return r.DB.Save(book).Error
}

func (r *BookRepository) Delete(book *database.Book) error {
	return r.DB.Delete(book).Error
}

func (r *BookRepository) FindByID(id uint) (*database.Book, error) {
	var book database.Book
	err := r.DB.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Find(author, genre string, start, end time.Time) ([]database.Book, error) {
	var books []database.Book

	db := r.DB

	if author != "" {
		db = db.Where("author ilike '%" + author + "%'")
	}

	if genre != "" {
		db = db.Where("genre ilike '%" + genre + "%'")
	}

	if start != (time.Time{}) && end != (time.Time{}) {
		db = db.Where("published_at BETWEEN ? AND ?", start, end)
	}

	err := db.Find(&books).Error
	if err != nil {
		return nil, err
	}

	return books, nil
}
