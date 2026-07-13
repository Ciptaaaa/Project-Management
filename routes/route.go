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
	 boardControl *controllers.BoardController){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error Loading .env file")
	}

	app.Post("/v1/auth/register", userControl.Register)
	app.Post("/v1/auth/login", userControl.Login)

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
}
