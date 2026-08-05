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

	//list
	listPosRepo := repositories.NewListPositionRepository()
	listRepo := repositories.NewListRepository()
	listService := services.NewListService(listRepo,boardRepo,listPosRepo)
	listController := controllers.NewListController(listService)
	
	//label
	labelRepo := repositories.NewLabelRepository()
	labelService := services.NewLabelService(labelRepo)
	labelController := controllers.NewLabelController(labelService)

	//card 
	cardRepo := repositories.NewCardRepository()
	cardService := services.NewCardService(cardRepo, listRepo, userRepo,labelRepo)
	cardController := controllers.NewCardController(cardService)

	//attachment
	attachmentRepo := repositories.NewAttachmentRepository(config.DB) 
	attachmentService := services.NewAttachmentService(attachmentRepo, cardRepo, userRepo)
	attachmentController := controllers.NewAttachmentController(attachmentService)

	//comment
	commentRepo := repositories.NewCommentRepository()
	commentService := services.NewCommentService(commentRepo, cardRepo, userRepo, listRepo, boardRepo) 
	commentController := controllers.NewCommentController(commentService)

	
	routes.Setup(app, userController, boardController, listController,cardController,labelController,attachmentController, commentController)

	port:= config.AppConfig.AppPort
	log.Println("Server is running on port:",port)

	log.Fatal(app.Listen(":"+port))
}