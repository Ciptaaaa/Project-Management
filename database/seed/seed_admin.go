package seed

import (
	"log"

	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/google/uuid"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("TiffanySayang")
	admin:= models.User{
		Name: "Super Admin",
		Email: "admin@example.com",
		Password: password,
		Role: "admin",
		PublicID: uuid.New(),
	}
	if err:= config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error;  err!=nil{
		log.Println("Failed to seed admin",err)
	}else{
		log.Println("Admin user seeded")
	}
}