package services

import (
	"errors"

	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Register(user *models.User) error
	Login(email,password string)(*models.User, error)
	GetByID(id uint)(*models.User, error)
	GetByPublicID(id string)(*models.User, error)
	GetAllPagination(filter, sort string, limit,ofset int) ([]models.User,int64, error )
	Update(user *models.User)error
	Delete(id uint) error
}

type userService struct{
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository)UserService{
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User) error {
//we checked email already regist (?)
existingUser, err := s.repo.FindByEmail(user.Email)


if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("database error")
	}

if existingUser != nil && existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}
//hash password

//set role admin
hashed, err := utils.HashPassword(user.Password)
if err != nil{
	return err
}

user.Password = hashed
user.Role = "user"
user.PublicID= uuid.New()
//save user
return s.repo.Create(user)
} 

func (s *userService)Login(email,password string)(*models.User, error){
user, err:= s.repo.FindByEmail(email)
if err != nil {
	return nil, errors.New("Invalid Credentials")
}
if !utils.CheckPasswordHash(password,user.Password){
return nil, errors.New("Invalid Credentials")
}
return user,nil
}

func (s *userService) GetByID(id uint)(*models.User, error){
return s.repo.FindByID(id)
}

func (s *userService) GetByPublicID(id string)(*models.User, error){
return s.repo.FindByPublicID(id)
}

func (s *userService) GetAllPagination(filter, sort string, limit,ofset int) ([]models.User,int64, error ){
	return s.repo.FindAllPagination(filter, sort, limit,ofset)
}

func (s *userService) Update(user *models.User)error {
	return s.repo.Update(user)
}

func (s *userService) Delete(id uint) error{
	return s.repo.Delete(id)
}