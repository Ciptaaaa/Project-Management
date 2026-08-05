package services

import (
	"errors"
	"time"

	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/google/uuid"
)

type CommentService interface {
	Create(cardPublicID, userPublicID, message string) (*models.Comment, error)
	GetByCardID(cardPublicID string) ([]models.Comment, error)
	Delete(publicID, requesterPublicID string) error 
}

type commentService struct {
	commentRepo repositories.CommentRepository
	cardRepo    repositories.CardRepository
	userRepo    repositories.UserRepository
	listRepo    repositories.ListRepository
	boardRepo   repositories.BoardRepository 
}

func NewCommentService(
	commentRepo repositories.CommentRepository,
	cardRepo repositories.CardRepository,
	userRepo repositories.UserRepository,
	listRepo repositories.ListRepository,
	boardRepo repositories.BoardRepository,
) CommentService {
	return &commentService{commentRepo, cardRepo, userRepo, listRepo, boardRepo}
}

func (s *commentService) Create(cardPublicID, userPublicID, message string) (*models.Comment, error) {
	if message == "" {
		return nil, errors.New("message is required")
	}
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return nil, errors.New("card not found")
	}
	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	comment := &models.Comment{
		PublicID:  uuid.New(),
		CardID:    card.InternalID,
		UserID:    user.InternalID,
		Message:   message,
		CreatedAt: time.Now(),
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	comment.User = models.UserLite{
		InternalID: user.InternalID,
		PublicID:   user.PublicID,
		Name:       user.Name,
		Email:      user.Email,
	}

	return comment, nil
}

func (s *commentService) GetByCardID(cardPublicID string) ([]models.Comment, error) {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return nil, errors.New("card not found")
	}
	return s.commentRepo.FindByCardID(card.InternalID)
}

func (s *commentService) Delete(publicID, requesterPublicID string) error {
	comment, err := s.commentRepo.FindByPublicID(publicID)
	if err != nil {
		return errors.New("comment not found")
	}

	requester, err := s.userRepo.FindByPublicID(requesterPublicID)
	if err != nil {
		return errors.New("requester not found")
	}

	isAuthor := requester.InternalID == comment.UserID
	isSiteAdmin := requester.Role == "admin"

	if isAuthor || isSiteAdmin {
		return s.commentRepo.Delete(uint(comment.InternalID))
	}

	isBoardOwner, err := s.isBoardOwnerOfCard(comment.CardID, requester.InternalID)
	if err != nil {
		return err
	}
	if !isBoardOwner {
		return errors.New("you are not allowed to delete this comment")
	}

	return s.commentRepo.Delete(uint(comment.InternalID))
}

func (s *commentService) isBoardOwnerOfCard(cardInternalID int64, requesterInternalID int64) (bool, error) {
	card, err := s.cardRepo.FindByID(uint(cardInternalID))
	if err != nil {
		return false, errors.New("card not found")
	}
	list, err := s.listRepo.FindByID(uint(card.ListID))
	if err != nil {
		return false, errors.New("list not found")
	}
	board, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		return false, errors.New("board not found")
	}
	return board.OwnerID == requesterInternalID, nil
}