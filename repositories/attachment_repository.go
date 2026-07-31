package repositories

import (
	"errors"

	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	GetByPublicId(publicId uuid.UUID)(*models.CardAttachment,error)
	FindByCardId(cardPublicID string)([]models.CardAttachment,error)
	Create(attachment *models.CardAttachment)error
	DeleteByPublicID(publicId uuid.UUID)error
}

type attachmentRepository struct{
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository{
	return &attachmentRepository{db}
}

func (r *attachmentRepository)GetByPublicId(publicId uuid.UUID)(*models.CardAttachment,error){
	var attachment models.CardAttachment
	if err := r.db.Where("public_id = ?",publicId).First(&attachment).Error; err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil,err
		}
	}
	return &attachment,nil
}

func (r *attachmentRepository) FindByCardId(cardPublicID string)([]models.CardAttachment,error){
	//get internal ID 
	var card models.Card
	if err := r.db.Where("public_id = ?",cardPublicID).First(&card).Error;err != nil{
		return nil, err
	}

	var attachments []models.CardAttachment
	if err := r.db.Where("card_internal_id = ?", card.InternalID).Find(&attachments).Error; err != nil{
		return nil,err
	}
	return attachments,nil
}

func (r *attachmentRepository) Create(attachment *models.CardAttachment)error{
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) DeleteByPublicID(publicId uuid.UUID)error{
	return r.db.Where("public_id = ?",publicId).Delete(&models.CardAttachment{}).Error; 
}
