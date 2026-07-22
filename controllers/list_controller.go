package controllers

import (
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}


func NewListController(s services.ListService) *ListController{
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx fiber.Ctx) error{
	list := new(models.List)
	if err := ctx.Bind().Body(list); err !=nil {
		return utils.BadRequest(ctx, "Failed Bind data ",err.Error())
	}
	if err := c.service.Create(list); err != nil{
		return utils.BadRequest(ctx, "Failed Create List",err.Error())
	}
	return utils.Success(ctx, "Successfully Created List!",list)
}

func (c *ListController) UpdateList(ctx fiber.Ctx) error{
	publicID := ctx.Params("id")
	list := new (models.List)

	if err := ctx.Bind().Body(list); err != nil {
		return utils.BadRequest(ctx, "Failed Parsed Data",err.Error())
	}
	if _,err:=  uuid.Parse(publicID); err != nil{
		return utils.BadRequest(ctx, "ID not valid",err.Error())
	}

existingList, err := c.service.GetByPublicID(publicID)
if err != nil {
	return utils.NotFound(ctx, "List not found!",err.Error())
}

list.InternalID =existingList.InternalID
list.PublicID = existingList.PublicID

if err := c.service.Update(list); err != nil {
	return utils.BadRequest(ctx, "Failed List Update",err.Error())
}

updatedList, err := c.service.GetByPublicID(publicID)
if err != nil { 
	return utils.NotFound(ctx, "Notfound List",err.Error())
}
return utils.Success(ctx, "Successfully Updated List",updatedList)
}

func (c *ListController) GetListOnBoard(ctx fiber.Ctx)error{
	boardPublicID := ctx.Params("board_id")
		if _,err:=  uuid.Parse(boardPublicID); err != nil{
		return utils.BadRequest(ctx, "ID not valid",err.Error())
	}
	lists, err :=c.service.GetByBoardID(boardPublicID)
	if err != nil{
		return utils.NotFound(ctx, "Failed get data!",err.Error())
	}
	return utils.Success(ctx, "Successfully Get Data",lists)
}