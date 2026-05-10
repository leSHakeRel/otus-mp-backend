package repositories

import (
	"movie-night-planner-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EveningFilmRepository struct {
	db *gorm.DB
}

func NewEveningFilmRepository(db *gorm.DB) *EveningFilmRepository {
	return &EveningFilmRepository{db: db}
}

func (r *EveningFilmRepository) Create(eveningFilm *models.EveningFilm) error {
	return r.db.Create(eveningFilm).Error
}

func (r *EveningFilmRepository) FindByID(id uuid.UUID) (*models.EveningFilm, error) {
	var eveningFilm models.EveningFilm
	err := r.db.First(&eveningFilm, id).Error
	if err != nil {
		return nil, err
	}
	return &eveningFilm, nil
}

func (r *EveningFilmRepository) FindByEveningID(eveningID uuid.UUID) ([]models.EveningFilm, error) {
	var eveningFilms []models.EveningFilm
	err := r.db.Where("evening_id = ?", eveningID).Find(&eveningFilms).Error
	if err != nil {
		return nil, err
	}
	return eveningFilms, nil
}

func (r *EveningFilmRepository) FindByTMDBID(tmdbID int) (*models.EveningFilm, error) {
	var eveningFilm models.EveningFilm
	err := r.db.Where("tmdb_id = ?", tmdbID).First(&eveningFilm).Error
	if err != nil {
		return nil, err
	}
	return &eveningFilm, nil
}

func (r *EveningFilmRepository) FindByEveningIDAndTMDBID(eveningID uuid.UUID, tmdbID int) (*models.EveningFilm, error) {
	var eveningFilm models.EveningFilm
	err := r.db.Where("evening_id = ? AND tmdb_id = ?", eveningID, tmdbID).First(&eveningFilm).Error
	if err != nil {
		return nil, err
	}
	return &eveningFilm, nil
}

func (r *EveningFilmRepository) Update(eveningFilm *models.EveningFilm) error {
	return r.db.Save(eveningFilm).Error
}

func (r *EveningFilmRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.EveningFilm{}, id).Error
}

func (r *EveningFilmRepository) DeleteByEveningID(eveningID uuid.UUID) error {
	return r.db.Where("evening_id = ?", eveningID).Delete(&models.EveningFilm{}).Error
}
