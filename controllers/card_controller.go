package controllers

import (
	"time"

	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/services"
	"github.com/adamabiyuu/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	service services.CardService
}

func NewCardController(s services.CardService) *CardController {
	return &CardController{service: s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
	}

	if err := c.service.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal Membuat card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dibuat", card)
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	type updateCardRequest struct {
		ListPublicID string     `json:"list_id"`
		Title        string     `json:"title"`
		Description  string     `json:"description"`
		DueDate      *time.Time `json:"due_date"`
		Position     int        `json:"position"`
	}

	var req updateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Id tidak valid", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Position:    req.Position,
		PublicID:    uuid.MustParse(publicID),
	}

	if err := c.service.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal update data", err.Error())
	}

	return utils.Success(ctx, "Card Berhasil diperbaharui", card)
}

func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Card tidak ditemukan", err.Error())
	}

	if err := c.service.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus data", err.Error())
	}
	return utils.Success(ctx, "Card berhasil dihapus", publicID)
}

func (c *CardController) GetListCard(ctx *fiber.Ctx) error {
	listID := ctx.Params("list_id")
	if _, err := uuid.Parse(listID); err != nil {
		return utils.BadRequest(ctx, "Id list tidak valid", err.Error())
	}

	cards, err := c.service.GetByListID(listID)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Mengambil data", err.Error())
	}

	return utils.Success(ctx, "Data Card Berhasil Diambil", cards)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")

	card, err := c.service.GetByPublicID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Error saat mengambil data", err.Error())
	}
	if card == nil {
		return utils.NotFound(ctx, "Card tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "Data berhasil diambil", card)

}

func (c *CardController) AddCardLabel(ctx *fiber.Ctx) error {
	cardId := ctx.Params("id")

	var body struct {
		LabelID string `json:"label_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Id tidak valid", err.Error())
	}
	if err := c.service.AddLabel(cardId, body.LabelID); err != nil {
		return utils.BadRequest(ctx, "Gagal menambahkan data label", err.Error())
	}
	return utils.Success(ctx, "Label Berhasil ditambahkan", nil)
}
func (c *CardController) RemoveCardLabel(ctx *fiber.Ctx) error {
	cardId := ctx.Params("id")

	var body struct {
		LabelID string `json:"label_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Id tidak valid", err.Error())
	}
	if err := c.service.RemoveLabel(cardId, body.LabelID); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus data label", err.Error())
	}
	return utils.Success(ctx, "Label Berhasil dihapus", nil)
}



// func (c *CardController) UploadAttachment(ctx *fiber.Ctx) error {

// 	cardID := ctx.Params("id")

// 	file, err := ctx.FormFile("file")
// 	if err != nil {
// 		return utils.BadRequest(ctx, "File wajib diupload", err.Error())
// 	}

// 	filename := uuid.New().String() + "_" + file.Filename

// 	if err := ctx.SaveFile(file, "./uploads/"+filename); err != nil {
// 		return utils.InternalServerError(ctx, "Gagal menyimpan file", err.Error())
// 	}

// 	// Ambil user dari JWT (sesuaikan dengan projectmu)
// 	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
// 	userID := claims["pub_id"].(string)

// 	attachment, err := c.attachmentService.Create(cardID, userID, filename)
// 	if err != nil {
// 		return utils.InternalServerError(ctx, "Upload gagal", err.Error())
// 	}

// 	return utils.Success(ctx, "Attachment berhasil diupload", attachment)
// }

// func (c *CardController) GetAttachments(ctx *fiber.Ctx) error {

// 	cardID := ctx.Params("id")

// 	attachments, err := c.attachmentService.FindByCardID(cardID)
// 	if err != nil {
// 		return utils.InternalServerError(ctx, "Gagal mengambil attachment", err.Error())
// 	}

// 	return utils.Success(ctx, "Data attachment berhasil diambil", attachments)
// }

// func (c *CardController) DeleteAttachment(ctx *fiber.Ctx) error {

// 	attachmentID := ctx.Params("attachment_id")

// 	pubID, err := uuid.Parse(attachmentID)
// 	if err != nil {
// 		return utils.BadRequest(ctx, "Attachment ID tidak valid", err.Error())
// 	}

// 	if err := c.attachmentService.DeleteByPublicID(pubID); err != nil {
// 		return utils.InternalServerError(ctx, "Gagal menghapus attachment", err.Error())
// 	}

// 	return utils.Success(ctx, "Attachment berhasil dihapus", nil)
// }