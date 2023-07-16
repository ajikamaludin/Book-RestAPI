package database

import "time"

type Book struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Author      string    `gorm:"not null"`
	PublishedAt time.Time `gorm:"not null" json:"published_at"`
	Edition     string
	Description string
	Genre       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Collection struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Books     []Book `gorm:"many2many:collection_books;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CollectionBook struct {
	CollectionID uint
	BookID       uint
}
