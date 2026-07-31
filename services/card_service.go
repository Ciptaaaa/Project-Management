package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/models/types"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint)error
	GetByListID(listPublicID string)([]models.Card,error)
	GetByID(id uint)(*models.Card,error)
	GetByPublicID(publicID string)(*models.Card,error)
	AddLabel(cardPublicID, labelPublicID string) error
	RemoveLabel(cardPublicID, labelPublicID string) error
	UpdatePositions(listPublicID string, cardOrder []string) error
}


type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
	labelRepo repositories.LabelRepository
}

func NewCardService (
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepository,
	userRepo repositories.UserRepository,
	labelRepo repositories.LabelRepository,
) CardService{
	return &cardService{cardRepo,listRepo,userRepo,labelRepo}
}
func sortCardByPosition(cards []models.Card, order []uuid.UUID)[]models.Card{
//map for fast searching

orderMap := make(map[uuid.UUID]int)
for i,id := range order{
	orderMap[id]= i
}

defaultIndex := len(order)

//sort slice 
sort.SliceStable( cards, func(i,j int)bool{
idx,ok :=  orderMap[cards[i].PublicID]
if !ok {
	idx = defaultIndex
}
idxJ,okJ := orderMap[cards[j].PublicID]
if !okJ {
	idxJ = defaultIndex
}

if idx == idxJ {
	return cards[i].CreatedAt.Before(cards[j].CreatedAt)
}
return idx < idxJ
})
return cards
}

func (s *cardService) Create(card *models.Card, listPublicID string) error{
	// Get List from listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w",err)
	}

	//set list_internal_id
	card.ListID= list.InternalID
	//Generate Public_id if not set
	if card.PublicID == uuid.Nil{
		card.PublicID = uuid.New()
	}
	card.CreatedAt= time.Now()

	//Transaction
	tx:= config.DB.Begin()
	defer func (){
	if r := recover(); r !=nil{
		tx.Rollback()
		panic(r)
	}
	}()

	//save card
	if err := tx.Create(card).Error; err !=nil{
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w ",err)
	}
	//update or create card position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).Where("list_internal_id = ?",list.InternalID).First(&position).Error; err != nil{
		if errors.Is(err , gorm.ErrRecordNotFound){
			//create card position if not exist
			position = models.CardPosition{
				PublicID: uuid.New(),
				ListID: list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}
			if err := tx.Create(&position).Error; err !=nil{
				tx.Rollback()
				return fmt.Errorf("failed to create card position: %w",err)
			}
		}else{
			tx.Rollback()
			return fmt.Errorf("failed to get card position: %w",err)
		}
	}else{
		//added new card to list 
		position.CardOrder = append(position.CardOrder, card.PublicID)
		if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?",position.InternalID).Update("card_order",position.CardOrder).Error; err != nil{
			tx.Rollback()
			return fmt.Errorf("failed to update card position : %w",err)
		}
	}
	//commit transaction
	if err:= tx.Commit().Error; err !=nil{
		return fmt.Errorf("Transaction commit failed: %w",err)
	}

	return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error{
	//get old card
	existingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found : %w",err)
	}
	//get new list
	newList,err := s.listRepo.FindByPublicID(listPublicID)
	if err !=nil{
		return fmt.Errorf("list not found: %w",err)
	}
	//transaction
	tx := config.DB.Begin()
	defer func (){
		if r := recover(); r !=nil{
			tx.Rollback()
			panic(r)
		}
	}()

	if existingCard.ListID != newList.InternalID {
		//delete old list
		var oldPos models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPos).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return fmt.Errorf("failed to get old card position: %w", err)
			}
		} else {
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}

			//update
			if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", oldPos.InternalID).Update("card_order",types.UUIDArray(filtered)).Error; err!=nil{
				tx.Rollback()
				return fmt.Errorf("Failed to update old card position: %w",err)
			}
		}
		//added new list
		var newPos models.CardPosition
		res:= tx.Where("list_internal_id = ?",newList.InternalID).First(&newPos)
		if errors.Is(res.Error, gorm.ErrRecordNotFound){
			newPos= models.CardPosition{
				PublicID: uuid.New(),
				ListID: newList.InternalID,
				CardOrder: types.UUIDArray{existingCard.PublicID},
			}
			if err := tx.Create(&newPos).Error;err!=nil{
				tx.Rollback()
				return fmt.Errorf("failed to create card position for new list: %w",err)
			}
		}else if res.Error == nil {
			updateOrder := append(newPos.CardOrder, existingCard.PublicID)
			if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?",newPos.InternalID).Update("card_order", types.UUIDArray(updateOrder)).Error;err != nil{
				tx.Rollback()
				return fmt.Errorf("failed to update new card position: %w",err)
			}
		} else{
			tx.Rollback()
			return fmt.Errorf("failed to get new card position: %w",res.Error)
		}
	}
	//update data card
	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListID = existingCard.ListID
	card.CreatedAt = existingCard.CreatedAt

	if err := tx.Save(card).Error; err !=nil{
		tx.Rollback()
		return fmt.Errorf("Failed to update card: %w",err)
	}
	//commit
	if err := tx.Commit().Error; err !=nil{
		tx.Rollback()
		return fmt.Errorf("Transaction commit failed: %w",err)
	}

	return nil

}
func (s *cardService) Delete(id uint)error{
return s.cardRepo.Delete(id)
}

func (s *cardService) GetByListID(listPublicID string)([]models.Card,error){
//verify list

list, err := s.listRepo.FindByPublicID(listPublicID)
if err !=nil{
	return nil, fmt.Errorf("List not found: %w",err)
}
//get card position
position, err := s.cardRepo.FindCardPositionByListID(list.InternalID)
if err != nil {
	return nil, fmt.Errorf("failed to get card position: %w",err)
}
//get all card in list
cards,err := s.cardRepo.FindByListID(listPublicID)
if err !=nil{
	return nil, fmt.Errorf("failed to get card: %w",err)
}
//sort

if position != nil && len(position.CardOrder)>0 {
	cards = sortCardByPosition(cards,position.CardOrder)
}
return cards,nil
}
func (s *cardService) GetByID(id uint)(*models.Card,error){
return s.cardRepo.FindByID(id)
}
func (s *cardService) GetByPublicID(publicID string)(*models.Card,error){
return s.cardRepo.FindByPublicID(publicID)
}

func (s *cardService) AddLabel(cardPublicID,labelPublicID string)error{
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil{
	return errors.New("Card Not found!")
	}
	label,err := s.labelRepo.FindByPublicID(labelPublicID)
	if err != nil {
		return errors.New("Label Not found!")
	}
	return s.cardRepo.AddLabel(uint(card.InternalID), uint(label.InternalID))
}

func (s *cardService) RemoveLabel(cardPublicID, labelPublicID string) error {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}
	label, err := s.labelRepo.FindByPublicID(labelPublicID)
	if err != nil {
		return errors.New("label not found")
	}
	return s.cardRepo.RemoveLabel(uint(card.InternalID), uint(label.InternalID))
}
func (s *cardService) UpdatePositions(listPublicID string, cardOrder []string) error {
	// verify the list exists before touching card_positions
	if _, err := s.listRepo.FindByPublicID(listPublicID); err != nil {
		return errors.New("list not found")
	}
	if err := s.cardRepo.UpdatePosition(listPublicID, cardOrder); err != nil {
		return fmt.Errorf("failed to update card position: %w", err)
	}
	return nil
}