package services

import (
	"errors"

	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/google/uuid"
)

type BoardService interface {

}

type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo repositories.UserRepository
}

func newBoardService(boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,) BoardService {
	return &boardService{boardRepo,userRepo }
}

func (s *boardService) Create (board *models.Board) error {
	user, err:= s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil{
		return errors.New("owner not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepo.Create(board)
}