package services

import (
	"errors"
	"fmt"

	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/models/types"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type listService struct {
	listRepo repositories.ListRepository
	boardRepo repositories.BoardRepository
	listPosRepo repositories.ListPositionRepository
}

type ListWithOrder struct{
	Positions []uuid.UUID
	Lists []models.List	
}


type ListService interface{
	GetByBoardID(boardPublicID string)(*ListWithOrder, error)
	GetById(id uint) (*models.List,error)
	GetByPublicID(publicID string)(*models.List, error)
	Create(list *models.List)error
	Update(list *models.List)error
	Delete(id uint)error
	UpdatePositions(boardPublicID string, positions []uuid.UUID)error
}



func NewListService (listRepo repositories.ListRepository,
	boardRepo repositories.BoardRepository,
	listPosRepo repositories.ListPositionRepository) ListService {
return &listService{listRepo,boardRepo,listPosRepo}
}


func ( s *listService) GetByBoardID(boardPublicID string)(*ListWithOrder, error) {
	//Verify Board
	_,err:=s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("Board not found!")
	}

	positions,err:=s.listPosRepo.GetListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("Failed to get list order :"+err.Error())
	}

	lists, err := s.listRepo.FindByBoardID(boardPublicID)
	if err != nil{
		return nil, errors.New("Failed to get list :"+err.Error()) 
	}

	//sorting by position
	orderedList := utils.SortListByPosition(lists,positions)

	return &ListWithOrder{
		Positions: positions,
		Lists: orderedList,
	},nil

}

func (s *listService) GetById(id uint) (*models.List,error){
	return s.listRepo.FindByID(id)
}

func (s *listService) GetByPublicID(publicID string)(*models.List, error){
	 return s.listRepo.FindByPublicID(publicID)
}

func (s *listService) Create(list *models.List)error {
	//validate Board
	board,err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return errors.New("board not found!")
		}
		return fmt.Errorf("Failed to get board : %w",err)
	} 
	list.BoardInternalID = board.InternalID
	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}
	//Transaction
	tx := config.DB.Begin()
	defer func (){
		if r := recover(); r !=nil {
			tx.Rollback()
		}
	}()


	if err := tx.Create(list).Error; err != nil{
		tx.Rollback()
		return fmt.Errorf("Failed to create list : %w",err)
	}

	//Update position
	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)
	if errors.Is(res.Error,gorm.ErrRecordNotFound){
		position = models.ListPosition{
			PublicID: uuid.New(),
			BoardID: board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}
			if err := tx.Create(&position).Error; err != nil{
		tx.Rollback()
		return fmt.Errorf("Failed to Create list position : %w",err)
	}
	}else if res.Error !=nil{
		tx.Rollback()
		return fmt.Errorf("Failed to Create list position : %w",res.Error)
	}else {
		position.ListOrder = append(position.ListOrder, list.PublicID)
		if err := tx.Model(&position).Update("list_order",position.ListOrder).Error; err != nil{
				tx.Rollback()
		return fmt.Errorf("Failed to Update list position : %w",err)
		}
	}
	//Commit trx
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Transaction commit failed : %w",err)
	}
	return nil
}

func (s *listService) Update(list *models.List)error{
	return s.listRepo.Update(list)
}

func (s *listService) Delete(id uint)error{
	return s.listRepo.Delete(id)
}

func (s *listService) UpdatePositions(boardPublicID string, positions []uuid.UUID)error{
	//verify Board
	board,err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil { 
		return errors.New("Board not found")
	}
	//get list position
	position,err := s.listPosRepo.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("List position not found")
	}
	//update list order
	position.ListOrder = positions
	return s.listPosRepo.UpdateListOrder(position)
}