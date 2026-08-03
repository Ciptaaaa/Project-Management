package controllers

import (
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AttachmentController struct {
	service services.AttachmentService
}

func NewAttachmentController(s services.AttachmentService) *AttachmentController {
	return &AttachmentController{service: s}
}

func (c *AttachmentController) UploadAttachments(ctx fiber.Ctx) error {
	cardID := ctx.Params("id")

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return utils.BadRequest(ctx, "Failed to read uploaded file", err.Error())
	}

	file, err := fileHeader.Open()
	if err != nil {
		return utils.BadRequest(ctx, "Failed to open uploaded file", err.Error())
	}
	defer file.Close()

	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return utils.BadRequest(ctx, "Invalid Token Claims", "Token Claims not valid")
	}
	userPublicID, ok := claims["public_id"].(string)
	if !ok {
		return utils.BadRequest(ctx, "Invalid Token Claims", "public_id claim missing")
	}

	attachment, err := c.service.Create(cardID, userPublicID, file, fileHeader.Filename)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to upload attachment", err.Error())
	}
	return utils.Success(ctx, "Successfully uploaded attachment", attachment)
}

func (c *AttachmentController) GetAttachments(ctx fiber.Ctx) error {
	cardID := ctx.Params("id")
	attachments, err := c.service.GetByCardID(cardID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to get attachments", err.Error())
	}
	return utils.Success(ctx, "Successfully get attachments", attachments)
}

func (c *AttachmentController) DeleteAttachments(ctx fiber.Ctx) error {
	attachmentID := ctx.Params("attachment_id")
	publicID, err := uuid.Parse(attachmentID)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid attachment id", err.Error())
	}

	if err := c.service.DeleteByPublicID(publicID); err != nil {
		return utils.InternalServerError(ctx, "Failed to delete attachment", err.Error())
	}
	return utils.Success(ctx, "Successfully deleted attachment", attachmentID)
}