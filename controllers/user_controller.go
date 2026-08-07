package controllers

import (
	"math"
	"strconv"
	"time"

	"github.com/Ciptaaaa/Project-Management.git/models"
	"github.com/Ciptaaaa/Project-Management.git/services"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
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

/*
	Sebelumnya token dikirim di body JSON, lalu frontend simpan di localStorage.
	Sekarang di-set sebagai HttpOnly cookie yang tidak bisa dibaca JavaScript,
	jadi celah XSS tidak bisa langsung curi session.

	Access token pendek (15 menit). Refresh token panjang tapi Path-nya dibatasi
	ke /v1/auth/refresh saja, jadi tidak ikut terkirim di sembarang request.
*/
utils.SetAuthCookies(ctx, token, refreshToken, 15*time.Minute, 7*24*time.Hour)

var userResponse models.UserResponse
_=  copier.Copy(&userResponse, &user)

/*
	Body tetap dikembalikan untuk backward-compatibility selama frontend masih
	transisi — setelah frontend berhenti membaca body, baris ini boleh dihapus.
	Cookie sudah di-set lewat utils.SetAuthCookies di atas.
*/
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

func (c *UserController) RefreshToken(ctx fiber.Ctx) error {
	/*
		Coba dari cookie dulu — ini yang dipakai frontend SPA setelah migrasi.
		Kalau kosong, baca dari body JSON sebagai backward-compatibility untuk
		client yang belum pakai cookie.
	*/
	refreshTokenStr := ctx.Cookies(utils.RefreshCookieName)

	if refreshTokenStr == "" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := ctx.Bind().Body(&body); err != nil {
			return utils.BadRequest(ctx, "Invalid Request", err.Error())
		}
		refreshTokenStr = body.RefreshToken
	}

	if refreshTokenStr == "" {
		return utils.BadRequest(ctx, "Invalid Request", "refresh_token is required")
	}

	userID, err := utils.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		return utils.Unauthorized(ctx, "Refresh failed", err.Error())
	}
	user, err := c.service.GetByID(uint(userID))
	if err != nil {
		return utils.Unauthorized(ctx, "Refresh failed", "user not found")
	}

	newAccessToken, err := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to generate token", err.Error())
	}

	/*
		Hanya access token yang diperbarui saat refresh, bukan refresh token.
		Refresh token tidak perlu diganti setiap kali — umurnya sendiri (7 hari)
		jauh lebih panjang dari access (15 menit).
		Parameter ketiga "" = tidak set refresh cookie baru.
	*/
	utils.SetAuthCookies(ctx, newAccessToken, "", 15*time.Minute, 0)

	/*
		Body dikembalikan untuk backward-compatibility selama ada client yang
		masih membaca token dari response JSON.
	*/
	return utils.Success(ctx, "Token refreshed successfully", fiber.Map{
		"access_token": newAccessToken,
	})
}

/*
Logout menghapus kedua cookie di sisi server.

Ini wajib ada begitu cookie jadi HttpOnly: frontend tidak punya cara menghapus
cookie yang tidak bisa dibacanya. Tanpa endpoint ini, tombol keluar hanya
membersihkan state React sementara browser tetap memegang session yang sah —
tekan reload, dan user "kembali login" tanpa memasukkan apa pun.

Sengaja tidak memvalidasi token dulu. Logout harus selalu berhasil; menolak
membersihkan cookie karena token sudah kedaluwarsa justru meninggalkan cookie
basi yang tidak bisa dibuang user lewat UI.
*/
func (c *UserController) Logout(ctx fiber.Ctx) error {
	utils.ClearAuthCookies(ctx)
	return utils.Success(ctx, "Logout Successfully!", nil)
}

/*
Me mengembalikan user pemilik cookie yang sedang dipakai.

Ini pengganti kebiasaan lama menyimpan objek user di localStorage. Nama, email,
dan role adalah data pribadi; menyimpannya di storage yang bisa dibaca skrip
mana pun tidak ada gunanya sekarang setelah tokennya sendiri sudah HttpOnly.
Frontend memanggil endpoint ini sekali saat halaman dibuka, menaruh hasilnya di
state React, dan membiarkannya hilang saat tab ditutup.

Claims sudah divalidasi middleware, jadi di sini tinggal dibaca. public_id
dipakai, bukan user_id, karena itu identitas yang dipakai seluruh API.
*/
func (c *UserController) Me(ctx fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		return utils.Unauthorized(ctx, "Invalid session", "claims missing")
	}

	publicID, ok := claims["public_id"].(string)
	if !ok {
		return utils.Unauthorized(ctx, "Invalid session", "public_id claim missing")
	}

	user, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.Unauthorized(ctx, "Invalid session", "user not found")
	}

	var userResp models.UserResponse
	if err := copier.Copy(&userResp, &user); err != nil {
		return utils.InternalServerError(ctx, "Error parsing data", err.Error())
	}

	return utils.Success(ctx, "Session valid", userResp)
}
