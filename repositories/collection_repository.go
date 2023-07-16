package repositories

import (
	"grace/database"

	"gorm.io/gorm"
)

type CollectionRepository struct {
	DB *gorm.DB
}

func (r *CollectionRepository) Create(collection *database.Collection) error {
	return r.DB.Create(collection).Error
}

func (r *CollectionRepository) Update(collection *database.Collection) error {
	return r.DB.Save(collection).Error
}

func (r *CollectionRepository) Delete(collection *database.Collection) error {
	return r.DB.Delete(collection).Error
}

func (r *CollectionRepository) FindByID(id uint) (*database.Collection, error) {
	var collection database.Collection
	err := r.DB.Preload("Books").First(&collection, id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}
