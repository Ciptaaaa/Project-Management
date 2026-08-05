package repositories

import (
	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByCardID(cardID int64) ([]models.Comment, error)
	FindByPublicID(publicID string) (*models.Comment, error)
	Delete(id uint) error
}

type commentRepository struct{}

func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return config.DB.Create(comment).Error
}

func (r *commentRepository) FindByCardID(cardID int64) ([]models.Comment, error) {
	var comments []models.Comment
	err := config.DB.
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("internal_id", "public_id", "name", "email")
		}).
		Where("card_internal_id = ?", cardID).
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) FindByPublicID(publicID string) (*models.Comment, error) {
	var comment models.Comment
	if err := config.DB.Where("public_id = ?", publicID).First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Comment{}, id).Error
}