package repositories

import (
	"github.com/adamabiyuu/project-management/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	FindByCardID(cardPublicID string) ([]models.CardAttachment, error)
	Create(attachment *models.CardAttachment) error
	DeleteByPublicID(publicID uuid.UUID) error
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db}
}

func (r *attachmentRepository) FindByCardID(cardPublicID string) ([]models.CardAttachment, error) {
	// ambil internal_id
	var card models.Card
	err := r.db.Where("card_public_id = ?", cardPublicID).First(&card).Error;
	if err != nil {
		return nil, err
	}
	var attachments []models.CardAttachment
	if err := r.db.Where("card_internal_id = ?", card.InternalID).Find(&attachments).Error; err != nil {
		return nil, err
	}
	
	return attachments, nil
}

func (r *attachmentRepository) Create(attachment *models.CardAttachment) error {
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) DeleteByPublicID(publicID uuid.UUID) error {
	return r.db.Where("public_id = ?", publicID).Delete(&models.CardAttachment{}).Error
}