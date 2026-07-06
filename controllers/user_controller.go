package controllers

import (
	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController{
	return &UserController{service:s}
}


func (c *UserController)Register (ctx fiber.Ctx)error{
	user := new(models.User)

	if err := ctx.Bind().Body(user); err != nil {
		return utils.BadRequest(ctx,"Failed Parsed Data", err.Error())
	}

	if err := c.service.Register(user); err != nil{
		return utils.BadRequest(ctx,"Failed Registration", err.Error())
	}
var userResponse models.UserResponse
_=  copier.Copy(&userResponse, &user)
	return utils.Success(ctx, "Register Success!", userResponse)
}

func (c *UserController) Login(ctx fiber.Ctx)error {
var body struct{
	Email string `json:"email"`
	Password string `json:"password"`
}
if err := ctx.Bind().Body(&body); err != nil{
	return utils.BadRequest(ctx, "Invalid Request",err.Error())
}
user,err:= c.service.Login(body.Email,body.Password)
if err != nil {
return utils.Unauthorized(ctx, "Login Failed!", err.Error())
}

token, err := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
if err != nil{
	return utils.BadRequest(ctx, "Failed to generate token", err.Error())
}
refreshToken, err := utils.GenerateRefreshToken(user.InternalID)
if err != nil{
	return utils.BadRequest(ctx, "Failed to generate refresh token",err.Error())
}
var userResponse models.UserResponse
_=  copier.Copy(&userResponse, &user)

return utils.Success(ctx, "Login Successfully!", fiber.Map{
	"access_token":token,
	"refresh_token":refreshToken,
	"user":userResponse,
})
}


func (c *UserController) GetUser(ctx fiber.Ctx) error{
id := ctx.Params("id")
user, err:= c.service.GetByPublicID(id)
if err != nil {
	return utils.NotFound(ctx, "Data not found!", err.Error())
}
var userResp models.UserResponse
err= copier.Copy(&userResp, &user)

if err != nil {
	return utils.BadRequest(ctx, "Internal Server Error:", err.Error())
}
return utils.Success(ctx, "Data Found!", userResp)
}