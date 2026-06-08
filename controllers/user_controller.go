// ctx hanya ada di layer Controller,
// karena tugasnya memang menjadi penghubung antara request/response HTTP dengan logika aplikasi (Service).
package controllers

import (
	"math"
	"strconv"

	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/services"
	"github.com/adamabiyuu/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Success(ctx, "Register Success", userResp)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "invalid request", err.Error())
	}
	
	user, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "invalid credential", err.Error())
	}

	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Success(ctx, "Login Successfully", fiber.Map{
		"access_token": token,
		"refresh_token": refreshToken,
		"user": userResp,
	})	
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "Data Not Found", err.Error())
	}
	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)
	if err != nil {
		return utils.BadRequest(ctx, "Internal Server Error", err.Error())
	}
	return utils.Success(ctx, "Data Berhasil Ditemukan", userResp)
}

func (c *UserController) GetUserPagination(ctx *fiber.Ctx) error {
	// /users/page?page=1&limit=10&sort=-id&filter=tryadi
	
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	users, total, err := c.service.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &users)

	meta := utils.PaginationMeta{
		Page: page,
		Limit: limit,
		Total: int(total),
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		Filter: filter,
		Sort: sort,
	}

	if total == 0 {
		return utils.NotFoundPagination(ctx, "Data pengguna Tidak ditemukan", userResp, meta)
	}

	return utils.SuccessPagination(ctx, "Data ditemukan", userResp, meta)
	
}


