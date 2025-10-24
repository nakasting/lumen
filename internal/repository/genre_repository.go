package repository

import (
	"lumen/internal/model"

	"gorm.io/gorm"
)

type GenreRepo interface {
	Create(genre *model.Genre) error
	FindAll() ([]model.Genre, error)
	Update(genre *model.Genre) error
	FindByID(id string) (*model.Genre, error)
	ExistsBySlug(slug string, id string) (bool, error)
	ExistsByID(id string) (bool, error)
}

type genreRepo struct {
	db *gorm.DB
}

func NewGenreRepo(db *gorm.DB) GenreRepo {
	return &genreRepo{db: db}
}

func (r *genreRepo) FindAll() ([]model.Genre, error) {
	var genres []model.Genre

	if err := r.db.Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

func (r *genreRepo) FindByID(id string) (*model.Genre, error) {
	var genre *model.Genre

	if err := r.db.First(&genre, id).Error; err != nil {
		return nil, err
	}
	return genre, nil
}

func (r *genreRepo) Create(genre *model.Genre) error {
	if err := r.db.Create(&genre).Error; err != nil {
		return err
	}
	return nil
}

func (r *genreRepo) Update(genre *model.Genre) error {
	if err := r.db.Save(&genre).Error; err != nil {
		return err
	}
	return nil
}

func (r *genreRepo) ExistsBySlug(slug string, id string) (bool, error) {
	var count int64

	if id == "" {
		if err := r.db.Model(&model.Genre{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
			return false, err
		}
	} else {
		if err := r.db.Model(&model.Genre{}).Where("slug = ? AND id <> ?", slug, id).Count(&count).Error; err != nil {
			return false, err
		}
	}
	return count > 0, nil
}

func (r *genreRepo) ExistsByID(id string) (bool, error) {
	var count int64

	if err := r.db.Model(&model.Genre{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
