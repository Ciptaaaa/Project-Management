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

	//user 
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController:= controllers.NewUserController(userService)

	//board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository() 
	boardService := services.NewBoardService(boardRepo, userRepo,boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	
	
	routes.Setup(app, userController, boardController)

	port:= config.AppConfig.AppPort
	log.Println("Server is running on port:",port)

	log.Fatal(app.Listen(":"+port))
}