package controllers

import (
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
)

type LabelController struct {
	service services.LabelService
}

func NewLabelController(s services.LabelService) *LabelController{
	return &LabelController{service:s}
}


func (c *LabelController) CreateLabel(ctx fiber.Ctx) error {
label := new(models.Label)
if err := ctx.Bind().Body(label); err != nil{
	return utils.BadRequest(ctx, "Failed parsed data",err.Error())
}
if label.Name == "" || label.Color == ""{
	return utils.BadRequest(ctx,"Failed create label","Name and color are required")
}

if err := c.service.Create(label);err != nil{
	return utils.BadRequest(ctx,"Failed save data", err.Error())
}
return utils.Success(ctx,"Successfully Created Label",label)
}

func (c *LabelController) GetLabelByID(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")
	label, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Label not found", err.Error())
	}
	return utils.Success(ctx, "Successfully Get Label", label)
}
func (c *LabelController) UpdateLabel(ctx fiber.Ctx) error{
	publicID := ctx.Params("id")
	existing, err := c.service.GetByPublicID(publicID)

	if err != nil { 
		return utils.NotFound(ctx, "Label not found",err.Error())
	}
	var body struct { 
		Name string `json:"name"`
		Color string `json:"color"`
	}

	if err := ctx.Bind().Body(&body); err != nil{
		return utils.BadRequest(ctx, "Failed Parsed Data",err.Error())
	}

	if body.Name != ""{
		existing.Name = body.Name
	}
	if body.Color != ""{
		existing.Color= body.Color
	}
	if err := c.service.Update(existing); err != nil{
		return utils.BadRequest(ctx, "Failed Updated Data",err.Error())
	}
	return utils.Success(ctx,"Successfully Updated Label", existing)
}
func (c *LabelController) DeleteLabel(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")
	label, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Label not found", err.Error())
	}

	if err := c.service.Delete(uint(label.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Failed Delete Data", err.Error())
	}
	return utils.Success(ctx, "Successfully Deleted Label", publicID)
}
func (c *LabelController) GetAllLabels(ctx fiber.Ctx) error {
	labels, err := c.service.GetAllLabel()
	if err != nil {
		return utils.InternalServerError(ctx, "Failed Get Data", err.Error())
	}
	return utils.Success(ctx, "Successfully Get All Label", labels)
}