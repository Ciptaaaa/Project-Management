package controllers

import (
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
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
