package controllers

import (
	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/services"
	"github.com/adamabiyuu/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membaca Request", err.Error())
	}
	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat List", err.Error())
	}

	return utils.Success(ctx, "List Berhasil Dibuat", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membaca Request", err.Error())
	}

	//validasi publicID benar format UUID atau tidak
	//misal 0000.aaa.00fasjlf.0000
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	//cek apakah list dengan publicID tersebut ada di database
	existingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}
	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Update List", err.Error())
	}

	//ambil data list terbaru
	// data list yang sudah terupdate
	updatedList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List Tidak Ditemukan", err.Error())
	}
	return utils.Success(ctx, "Berhasil Memperbaharui List", updatedList)
}