package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	InternalID int64     `json:"internal_id" gorm:"primaryKey;autoIncrement;column:internal_id"`
	PublicID   uuid.UUID `json:"public_id" gorm:"type:uuid;column:public_id"`
	CardID     int64     `json:"card_internal_id" gorm:"column:card_internal_id"`
	UserID     int64     `json:"user_internal_id" gorm:"column:user_internal_id"`
	Message    string    `json:"message" gorm:"column:message"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	User 	   UserLite  `json:"user" gorm:"foreignKey:UserID;references:InternalID"`
}

func (Comment) TableName() string {
	return "comments"
}