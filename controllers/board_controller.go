package controllers

import (
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController{
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard (ctx fiber.Ctx)error{
var userID uuid.UUID
var err error
	board:= new(models.Board)
	

	if err := ctx.Bind().Body(board); err !=nil{
		return utils.BadRequest(ctx, "Failed Request", err.Error())
	}
	
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	
	if !ok {
		return utils.BadRequest(ctx,"Invalid Token Claims","Token Claims not valid")
	}

	userID, err = uuid.Parse(claims["public_id"].(string))
	if err != nil{
	return utils.BadRequest(ctx, "Failed Request",err.Error())
	}

	board.OwnerPublicID=userID
	
	if err := c.service.Create(board); err != nil{
		return utils.BadRequest(ctx, "Failed Save Data", err.Error())
	}

	return utils.Success(ctx,"Successfully Created Board",board)
}