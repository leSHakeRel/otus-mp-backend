package repositories

import (
	"movie-night-planner-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByID(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) FindByEveningID(eveningID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("evening_id = ?", eveningID).Preload("User").Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) FindByUserID(userID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("user_id = ?", userID).Preload("Evening").Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *CommentRepository) DeleteByEveningID(eveningID uuid.UUID) error {
	return r.db.Where("evening_id = ?", eveningID).Delete(&models.Comment{}).Error
}

func (r *CommentRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Comment{}).Error
}
