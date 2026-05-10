package repositories

import (
	"movie-night-planner-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VoteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) Create(vote *models.Vote) error {
	return r.db.Create(vote).Error
}

func (r *VoteRepository) FindByID(id uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.Preload("User").Preload("EveningFilm").First(&vote, id).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) FindByEveningID(eveningID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("evening_id = ?", eveningID).Preload("User").Preload("EveningFilm").Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (r *VoteRepository) FindByEveningIDAndFilmID(eveningID uuid.UUID, eveningFilmID uuid.UUID) ([]models.Vote, error) {
	var votes []models.Vote
	err := r.db.Where("evening_id = ? AND evening_film_id = ?", eveningID, eveningFilmID).Preload("User").Find(&votes).Error
	if err != nil {
		return nil, err
	}
	return votes, nil
}

func (r *VoteRepository) FindByUserIDAndEveningID(userID uuid.UUID, eveningID uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.Where("user_id = ? AND evening_id = ?", userID, eveningID).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) FindByUserIDEveningIDAndFilmID(userID uuid.UUID, eveningID uuid.UUID, eveningFilmID uuid.UUID) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.Where("user_id = ? AND evening_id = ? AND evening_film_id = ?", userID, eveningID, eveningFilmID).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) Update(vote *models.Vote) error {
	return r.db.Save(vote).Error
}

func (r *VoteRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Vote{}, id).Error
}

func (r *VoteRepository) DeleteByEveningID(eveningID uuid.UUID) error {
	return r.db.Where("evening_id = ?", eveningID).Delete(&models.Vote{}).Error
}

func (r *VoteRepository) DeleteByUserIDAndEveningID(userID uuid.UUID, eveningID uuid.UUID) error {
	return r.db.Where("user_id = ? AND evening_id = ?", userID, eveningID).Delete(&models.Vote{}).Error
}
