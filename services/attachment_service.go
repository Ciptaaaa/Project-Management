package services

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type AttachmentService interface {
	GetByPublicID(publicId uuid.UUID) (*models.CardAttachment, error)
	GetByCardID(cardPublicID string) ([]models.CardAttachment, error)
	Create(cardPublicID, userPublicID string, file io.Reader, filename string) (*models.CardAttachment, error)
	DeleteByPublicID(publicId uuid.UUID) error
}

type attachmentService struct {
	attachmentRepo repositories.AttachmentRepository
	cardRepo       repositories.CardRepository
	userRepo       repositories.UserRepository
	cld            *cloudinary.Cloudinary
}

func NewAttachmentService(
	attachmentRepo repositories.AttachmentRepository,
	cardRepo repositories.CardRepository,
	userRepo repositories.UserRepository,
) AttachmentService {
	cld, err := cloudinary.NewFromURL(config.AppConfig.CloudinaryURL)
	if err != nil {
		// Handle the error appropriately
		panic("failed to initialize cloudinary client: " + err.Error())
	}
	return &attachmentService{attachmentRepo, cardRepo, userRepo, cld}
}

func (s *attachmentService) GetByPublicID(publicId uuid.UUID) (*models.CardAttachment, error) {
	return s.attachmentRepo.GetByPublicId(publicId)
}

func (s *attachmentService) GetByCardID(cardPublicID string) ([]models.CardAttachment, error) {
	return s.attachmentRepo.FindByCardId(cardPublicID)
}

func (s *attachmentService) Create(cardPublicID, userPublicID string, file io.Reader, filename string) (*models.CardAttachment, error) {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return nil, errors.New("card not found")
	}
	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	uploadResult, err := s.cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: "card-attachments",
	})
	if err != nil {
		return nil, errors.New("failed to upload file: " + err.Error())
	}

	attachment := &models.CardAttachment{
		PublicID:      uuid.New(),
		CardID:        card.InternalID,
		UserID:        user.InternalID,
		File:          uploadResult.SecureURL,
		CloudPublicID: uploadResult.PublicID,
		CreatedAt:     time.Now(),
	}

	if err := s.attachmentRepo.Create(attachment); err != nil {
		return nil, err
	}
	return attachment, nil
}

func (s *attachmentService) DeleteByPublicID(publicId uuid.UUID) error {
	attachment, err := s.attachmentRepo.GetByPublicId(publicId)
	if err != nil {
		return errors.New("attachment not found")
	}
	if _, err := s.cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: attachment.CloudPublicID,
	}); err != nil {
		return errors.New("failed to delete file from cloudinary: " + err.Error())
	}

	return s.attachmentRepo.DeleteByPublicID(publicId)
}