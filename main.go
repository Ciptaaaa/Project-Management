package main

import (
	"log"
	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/controllers"
	"github.com/Ciptaaaa/Project-Management.git/database/seed"
	"github.com/Ciptaaaa/Project-Management.git/repositories"
	"github.com/Ciptaaaa/Project-Management.git/routes"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/gofiber/fiber/v3"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	seed.SeedAdmin()
	app:= fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController:= controllers.NewUserController(userService)
	
	
	routes.Setup(app, userController)

	port:= config.AppConfig.AppPort
	log.Println("Server is running on port:",port)

	log.Fatal(app.Listen(":"+port))
}