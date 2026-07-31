package repositories

import (
	"fmt"
	"path/filepath"

	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/models/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardRepository interface {
	Create(card *models.Card) error
	Update(card *models.Card) (*models.Card, error) 
	Delete(id uint) error
	FindByID(id uint)(*models.Card, error)
	FindByPublicID(publicID string)(*models.Card, error)
	FindByListID(listID string) ([]models.Card,error)
	FindCardPositionByListID(id int64)(*models.CardPosition, error)
	UpdatePosition(listID string, position []string)error
	AddLabel(cardID uint, labelID uint) error
	RemoveLabel(cardID uint, labelID uint) error

}

type cardRepository struct {
}

func NewCardRepository() CardRepository{
	return &cardRepository{}
}


func (r *cardRepository) Create(card *models.Card)error{
return config.DB.Create(card).Error	
}
func (r *cardRepository) Update(card *models.Card) (*models.Card, error) {
	var existing models.Card
	if err := config.DB.Where("public_id = ?", card.PublicID).First(&existing).Error; err != nil {
		return nil, err
	}
	existing.Title = card.Title
	existing.Description = card.Description
	existing.DueDate = card.DueDate
	existing.Position = card.Position
	if err := config.DB.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}
func(r *cardRepository) Delete(id uint)error{
	return config.DB.Delete(&models.Card{},id).Error
}
func (r *cardRepository)FindByID(id uint)(*models.Card, error){
	var card models.Card
	err := config.DB.Preload("Labels").Preload("Assigness").First(&card,id).Error
	return &card,err
}

func (r *cardRepository) FindByPublicID(publicID string)(*models.Card, error){
	var card models.Card
	if err := config.DB.Preload("Assignees.User",func(tx *gorm.DB) *gorm.DB{
		return tx.Select("internal_id","public_id","name","email")
	}).Preload("Attachments").Where("public_id = ?",publicID).First(&card).Error; err!=nil{
		return nil,err
	}
	baseUrl := config.AppConfig.APPURL

	for i := range card.Attachments{
		card.Attachments[i].FileURL = fmt.Sprintf("%s/files/%s ",baseUrl,
		filepath.Base(card.Attachments[i].File),
	)
	}
	return &card,nil
}

func (r *cardRepository) FindByListID(listID string) ([]models.Card,error){
	var cards []models.Card
	err := config.DB.Joins("JOIN lists ON lists.internal_id = cards.list_internal_id").
	Where("lists.public_id = ?", listID).Order("position ASC").Find(&cards).Error
	return cards,err
}
func (r *cardRepository) FindCardPositionByListID(id int64)(*models.CardPosition,error){
	var position models.CardPosition
	err := config.DB.Where("list_internal_id = ?",id).First(&position).Error
	if err != nil{
		return nil, err
	}
	return &position,nil
}

func (r *cardRepository) UpdatePosition(listID string, position []string) error {
	uuids := make(types.UUIDArray, 0, len(position))
	for _, p := range position {
		u, err := uuid.Parse(p)
		if err != nil {
			return fmt.Errorf("invalid card id %q: %w", p, err)
		}
		uuids = append(uuids, u)
	}
	return config.DB.Model(&models.CardPosition{}).
		Where("list_internal_id = (SELECT internal_id FROM lists where public_id = ?)", listID).
		Update("card_order", uuids).Error
}
func (r *cardRepository) AddLabel(cardID uint, labelID uint) error {
	cardLabel := models.CardLabel{
		CardID:  int64(cardID),
		LabelID: int64(labelID),
	}
	return config.DB.Create(&cardLabel).Error
}

func (r *cardRepository) RemoveLabel(cardID uint, labelID uint) error {
	return config.DB.Where("card_internal_id = ? AND label_internal_id = ?", cardID, labelID).
		Delete(&models.CardLabel{}).Error
}