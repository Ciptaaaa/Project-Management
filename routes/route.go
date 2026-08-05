package routes

import (
	"log"

	"github.com/Ciptaaaa/Project-Management.git/controllers"
	"github.com/Ciptaaaa/Project-Management.git/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App,
	 userControl *controllers.UserController,
	 boardControl *controllers.BoardController,
	 listControl *controllers.ListController,
	 controllerControl *controllers.CardController,
	 labelControl *controllers.LabelController,
	 attachmentControl *controllers.AttachmentController,
	 commentControl *controllers.CommentController,
	 ){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error Loading .env file")
	}

	app.Post("/v1/auth/register", userControl.Register)
	app.Post("/v1/auth/login", userControl.Login)
	app.Post("/v1/auth/refresh", userControl.RefreshToken)

	//JWT PROTECTED ROUTES v2 
// 	api:= app.Group("/api/v1", jwtware.New(jwtware.Config{
// 		SigningKey: jwtware.SigningKey{
// 			Key: []byte(config.AppConfig.JWTSecret),
// 		},
// 		ContextKey:"user",
// 	  ErrorHandler: func(ctx *fiber.Ctx, err error) error {
//     return utils.Unauthorized(ctx, "Error Unauthorized!", err.Error())
// },
// 	}))

//go fiber v3 harus pake middleware 
	api:= app.Group("/api/v1", middleware.JWTProtected())

	//user group
	userGroup:= api.Group("/users")
	userGroup.Get("/page", userControl.GetUserPagination) // api/v1/users/page
	userGroup.Get("/:id", userControl.GetUser) // api/v1/users/:id
	userGroup.Put("/:id", userControl.UpdateUser)// update user
	userGroup.Delete("/:id",userControl.DeleteUser)//delete user soft deleted 


	//board group
	boardGroup := api.Group("/boards")
	boardGroup.Post("/",boardControl.CreateBoard)
	boardGroup.Put("/:id",boardControl.UpdateBoard)
	boardGroup.Post("/:id/members", boardControl.AddBoardMembers)
	boardGroup.Delete("/:id/members", boardControl.RemoveBoardMembers)
	boardGroup.Get("/my",boardControl.GetMyBoardPaginate)
	boardGroup.Get("/:board_id/lists",listControl.GetListOnBoard)
	boardGroup.Put("/:board_id/position",listControl.UpdateListPosition)

	//list group
	listGroup := api.Group("/lists")
	listGroup.Post("/",listControl.CreateList)
	listGroup.Put("/:id",listControl.UpdateList)
	listGroup.Delete("/:id",listControl.DeleteList)
	listGroup.Get("/:list_id/cards",controllerControl.GetListCard)
	listGroup.Put("/:list_id/positions",controllerControl.UpdatePosition) 	

	//Card group
	cardGroup:= api.Group("/cards")
	cardGroup.Post("/",controllerControl.CreateCard)
	cardGroup.Put("/:id",controllerControl.UpdateCard)
	cardGroup.Delete("/:id",controllerControl.DeleteCard)
	cardGroup.Get("/:id",controllerControl.GetCardDetail)
	cardGroup.Post("/:id/labels", controllerControl.AddCardLabel)
	cardGroup.Delete("/:id/labels", controllerControl.RemoveCardLabel)

	//Assign user to card
	cardGroup.Post("/:id/assignees", controllerControl.AssignCardUser)
	cardGroup.Delete("/:id/assignees", controllerControl.UnassignCardUser)
		
	//Attachment group
	cardGroup.Post("/:id/attachments", attachmentControl.UploadAttachments)
	cardGroup.Get("/:id/attachments", attachmentControl.GetAttachments)
	cardGroup.Delete("/:card_id/attachments/:attachment_id", attachmentControl.DeleteAttachments)

	//coment group
	cardGroup.Post("/:id/comments", commentControl.CreateComment)
	cardGroup.Get("/:id/comments", commentControl.GetComments)
	cardGroup.Delete("/:id/comments/:comment_id", commentControl.DeleteComment)

	//label group 
	labelGroup := api.Group("/labels")
	labelGroup.Post("/",labelControl.CreateLabel)
	labelGroup.Get("/",labelControl.GetAllLabels)
	labelGroup.Get("/:id",labelControl.GetLabelByID)
	labelGroup.Put("/:id",labelControl.UpdateLabel)
	labelGroup.Delete("/:id",labelControl.DeleteLabel)

	

	//Upload File
	// cardGroup.Post("/:id/attachments",controllerControl.UploadAttachments)//upload
	// cardGroup.Get("/:id/attachments",controllerControl.GetAttachments)//Get list attachments
	// cardGroup.Delete("/:card_id/attachments/:attachment_id",controllerControl.DeleteAttachments)//delete attachment
}
