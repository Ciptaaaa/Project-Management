package controllers

import (
	"math"
	"strconv"

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

func (c *BoardController) UpdateBoard(ctx fiber.Ctx)error{
	publicID := ctx.Params("id")
	board:= new (models.Board)

	if err := ctx.Bind().Body(board);err != nil{
		return utils.BadRequest(ctx, "Failed Parsed Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil{
		return utils.BadRequest(ctx, "ID not valid", err.Error())
	}
	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board Not Found!",err.Error())
	}
	board.InternalID= existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerID = existingBoard.OwnerID
	board.OwnerPublicID= existingBoard.OwnerPublicID
	board.CreatedAt= existingBoard.CreatedAt

	if err := c.service.Update(board); err != nil{
		return utils.BadRequest(ctx, "Faile Update Board", err.Error())
	}

	return utils.Success(ctx, "Successfully Update Board", board)
}

func (c *BoardController) AddBoardMembers(ctx fiber.Ctx) error{
	publicID := ctx.Params("id")

	var userIDs []string

	if err := ctx.Bind().Body(&userIDs);err != nil {
		return utils.BadRequest(ctx, "Failed Parsed Data",err.Error())
	}

	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		 return utils.BadRequest(ctx, "Failed Added Members",err.Error())
	}

	return utils.Success(ctx, "Successfully added member", nil)
}

func (c *BoardController) RemoveBoardMembers(ctx fiber.Ctx) error{
	publicID := ctx.Params("id")

	var userIDs []string

	if err := ctx.Bind().Body(&userIDs);err != nil {
		return utils.BadRequest(ctx, "Failed Parsed Data",err.Error())
	}

	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		 return utils.BadRequest(ctx, "Failed deleted Members",err.Error())
	}

	return utils.Success(ctx, "Successfully remove members", nil)
}

func (c *BoardController) GetMyBoardPaginate(ctx fiber.Ctx) error{
	claims,ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return utils.Unauthorized(ctx, "Invalid token claims", "")
	}

	userID, ok := claims["public_id"].(string)

	if !ok {
		return utils.Unauthorized(ctx, "Invalid token claims", "public_id missing or invalid")
	}

	page, _ := strconv.Atoi(ctx.Query("page","1"))
	limit, _ := strconv.Atoi(ctx.Query("limit","10"))
	offset := (page- 1) * limit

	filter := ctx.Query("filter","")
	sort:= ctx.Query("sort", "")

	boards,total, err := c.service.GetAllByUserPaginate(userID,filter,sort, limit,offset)
	if err != nil{
		return utils.InternalServerError(ctx, "Internal Server Error",err.Error())
	}
	meta := utils.PaginationMeta{
		Page:page,
		Limit: limit,
		Total: int(total),
		TotalPage: int(math.Ceil(float64(total)/float64(limit))),
		Filter: filter,
		Sort: sort,
	}

	return utils.SuccessPagination(ctx, "Successfully Get Data Board", boards, meta)
}