package services

import (
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
)

type LabelService interface {
	Create(label *models.Label)error
	Update(label *models.Label)error
	Delete(id uint)error
	GetByPublicID(publicID string)(*models.Label,error)
}

type labelService struct{
	repo repositories.LabelRepository
}
func NewLabelService(repo repositories.LabelRepository) LabelService{
	return &labelService{repo}
}

func (s *labelService) Create(label *models.Label) error {
	return s.repo.Create(label)
}
func (s *labelService) Update(label *models.Label) error {
	return s.repo.Update(label)
}
func (s *labelService) Delete(id uint) error {
	return s.repo.Delete(id)
}
func (s *labelService) GetByPublicID(publicID string) (*models.Label, error) {
	return s.repo.FindByPublicID(publicID)
}