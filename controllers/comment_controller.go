package controllers

import (
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type CommentController struct {
	service services.CommentService
}

func NewCommentController(s services.CommentService) *CommentController {
	return &CommentController{service: s}
}

func (c *CommentController) CreateComment(ctx fiber.Ctx) error {
	cardID := ctx.Params("id")

	var body struct {
		Message string `json:"message"`
	}
	if err := ctx.Bind().Body(&body); err != nil {
		return utils.BadRequest(ctx, "Failed Parsed Data", err.Error())
	}
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return utils.BadRequest(ctx, "Invalid Token Claims", "Token Claims not valid")
	}
	userPublicID, ok := claims["public_id"].(string)
	if !ok {
		return utils.BadRequest(ctx, "Invalid Token Claims", "public_id claim missing")
	}

	comment, err := c.service.Create(cardID, userPublicID, body.Message)
	if err != nil {
		return utils.BadRequest(ctx, "Failed create comment", err.Error())
	}
	return utils.Success(ctx, "Successfully created comment", comment)
}

func (c *CommentController) GetComments(ctx fiber.Ctx) error {
	cardID := ctx.Params("id")
	comments, err := c.service.GetByCardID(cardID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed get comments", err.Error())
	}
	return utils.Success(ctx, "Successfully get comments", comments)
}

func (c *CommentController) DeleteComment(ctx fiber.Ctx) error {
	commentID := ctx.Params("comment_id")

	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return utils.BadRequest(ctx, "Invalid Token Claims", "Token Claims not valid")
	}
	requesterPublicID, ok := claims["public_id"].(string)
	if !ok {
		return utils.BadRequest(ctx, "Invalid Token Claims", "public_id claim missing")
	}

	if err := c.service.Delete(commentID, requesterPublicID); err != nil {
		return utils.BadRequest(ctx, "Failed delete comment", err.Error())
	}
	return utils.Success(ctx, "Successfully deleted comment", commentID)
}