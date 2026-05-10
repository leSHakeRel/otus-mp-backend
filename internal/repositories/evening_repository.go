package repositories

import (
	"movie-night-planner-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EveningRepository struct {
	db *gorm.DB
}

func NewEveningRepository(db *gorm.DB) *EveningRepository {
	return &EveningRepository{db: db}
}

func (r *EveningRepository) Create(evening *models.Evening) error {
	return r.db.Create(evening).Error
}

func (r *EveningRepository) FindByID(id uuid.UUID) (*models.Evening, error) {
	var evening models.Evening
	err := r.db.Preload("Owner").Preload("EveningFilms").Preload("Comments.User").First(&evening, id).Error
	if err != nil {
		return nil, err
	}
	return &evening, nil
}

func (r *EveningRepository) FindByOwnerID(ownerID uuid.UUID, page, limit int) ([]models.Evening, int64, error) {
	var evenings []models.Evening
	var total int64

	r.db.Model(&models.Evening{}).Where("owner_id = ?", ownerID).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Preload("Owner").Preload("EveningFilms").Where("owner_id = ?", ownerID).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&evenings).Error
	if err != nil {
		return nil, 0, err
	}

	return evenings, total, nil
}

func (r *EveningRepository) FindAll(page, limit int, isPrivate *bool) ([]models.Evening, int64, error) {
	var evenings []models.Evening
	var total int64

	query := r.db.Model(&models.Evening{}).Preload("Owner").Preload("EveningFilms")

	if isPrivate != nil {
		query = query.Where("is_private = ?", *isPrivate)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&evenings).Error
	if err != nil {
		return nil, 0, err
	}

	return evenings, total, nil
}

func (r *EveningRepository) Update(evening *models.Evening) error {
	return r.db.Save(evening).Error
}

func (r *EveningRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Evening{}, id).Error
}
