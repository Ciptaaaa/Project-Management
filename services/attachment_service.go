package services

import (
	"errors"
	"time"

	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/google/uuid"
)

type AttachmentService interface {
	GetByPublicID(publicId uuid.UUID) (*models.CardAttachment, error)
	Create(CardPublicId,userPublicID,filename string)(*models.CardAttachment,error)
	DeleteByPublicID(publicId uuid.UUID)error
}

type attachmentService struct {
	attachmentRepo repositories.AttachmentRepository
	cardRepo repositories.CardRepository
	userRepo repositories.UserRepository
}

func NewAttachmentService(
	attachmentRepo repositories.AttachmentRepository,
	cardRepo repositories.CardRepository,
	userRepo repositories.UserRepository,
) AttachmentService{
	return &attachmentService{attachmentRepo,cardRepo,userRepo}
}

func (s *attachmentService) GetByPublicID(publicId uuid.UUID) (*models.CardAttachment, error){
	return s.attachmentRepo.GetByPublicId(publicId)
}

func (s *attachmentService) Create(CardPublicId,userPublicID,filename string)(*models.CardAttachment,error){
	card, err := s.cardRepo.FindByPublicID(CardPublicId)
	if err != nil {
		return nil,errors.New("Card not found!")
	}
	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil{
		return nil,errors.New("User not found!")
	}
	attachment := &models.CardAttachment{
		PublicID: uuid.New(),
		CardID: card.InternalID,
		UserID: user.InternalID,
		File: filename,
		CreatedAt: time.Now(),
	}

	if err := s.attachmentRepo.Create(attachment); err != nil{
		return nil,err
	}
	return attachment,nil
}

func (s *attachmentService) DeleteByPublicID(publicId uuid.UUID)error{
	return s.attachmentRepo.DeleteByPublicID(publicId)
}