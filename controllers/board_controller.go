package controllers

import (
	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/services"
	"github.com/adamabiyuu/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error

	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Membaca Request", err.Error())
	}

	userID, err = uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Gagal Membaca Request", err.Error())
	}
	board.OwnerPublicId = userID
	
	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Menyimpan Data", err.Error())
	}
	return utils.Success(ctx, "Board Berhasil Dibuat", board)
}

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID Tidak Valid", err.Error())
	}
	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board Tidak Ditemukan", err.Error())
	}
	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerID = existingBoard.OwnerID
	board.OwnerPublicId = existingBoard.OwnerPublicId
	board.CreatedAt = existingBoard.CreatedAt

	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Board", err.Error())
	}
	return utils.Success(ctx, "Board Berhasil Diperbaharui", board)
}

func (c *BoardController) AddBoardMembers(ctx *fiber.Ctx) error{
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}
	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Members", err.Error())
	}
	return utils.Success(ctx, "Member Berhasil ditambahkan", nil)
}