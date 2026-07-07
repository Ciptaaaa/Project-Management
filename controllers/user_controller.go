package controllers

import (
	"math"
	"strconv"

	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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



const (
    defaultLimit = 10
    maxLimit     = 100
    defaultPage  = 1
)

func (c *UserController) GetUserPagination (ctx fiber.Ctx) error {
	 page, err := strconv.Atoi(ctx.Query("page", "1"))
    if err != nil || page < 1 {
        page = defaultPage
    }
	 limit, err := strconv.Atoi(ctx.Query("limit", "10"))
    if err != nil || limit < 1 {
        limit = defaultLimit
    }
    if limit > maxLimit {
        limit = maxLimit 
    }
	offset := (page - 1 )* limit
	filter := ctx.Query("filter","")
	sort := ctx.Query("sort","")


	users,total,err := c.service.GetAllPagination(filter,sort,limit,offset,)
	if err != nil{
		return utils.BadRequest(ctx, "Failed Get Data", err.Error())
	}

	var userResp []models.UserResponse
	if err := copier.Copy(&userResp, &users); err != nil {
    return utils.BadRequest(ctx, "Failed to process data", err.Error())
	}

	
	meta := utils.PaginationMeta{
		Page:page,
		Limit:limit,
		Total: int(total),
		TotalPage:int (math.Ceil(float64(total)/(float64(limit)))),
		Filter: filter,
		Sort: sort,
	}

	if total == 0 {
		return utils.NotFoundPagination(ctx, "Data not found", userResp,meta)
	}

	return utils.SuccessPagination(ctx, "Data found",userResp, meta)
}


func (c *UserController) UpdateUser (ctx fiber.Ctx) error {
id := ctx.Params("id")
publicID, err := uuid.Parse(id)
if err != nil{
	return utils.BadRequest(ctx, "Invalid ID Format",err.Error())
}
var user models.User

if err := ctx.Bind().Body(&user); err != nil{
	return utils.BadRequest(ctx, "Failed Parsing Data",err.Error())
}

user.PublicID= publicID

if err := c.service.Update(&user);err != nil{
	return utils.BadRequest(ctx, "Failed Update Data",err.Error())
}

userUpdated, err := c.service.GetByPublicID(id)

if err != nil {
	return utils.InternalServerError(ctx, "Failed Receive Data",err.Error())
}

var userResp models.UserResponse
err= copier.Copy(&userResp, userUpdated)

if err != nil { 
	return utils.InternalServerError(ctx, "Error parsing data",err.Error())
}

return utils.Success(ctx, "Successfully Updated Data",userResp)
}

func (c *UserController) DeleteUser(ctx fiber.Ctx)error {
id,_:= strconv.Atoi(ctx.Params("id"))
if err := c.service.Delete(uint(id));err!=nil{
	return utils.InternalServerError(ctx, "Internal Server Error!", err.Error())
}
return utils.Success(ctx, "Successfully Delete user",id)
}
