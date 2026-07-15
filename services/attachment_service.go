package services

import (
	"errors"
	"time"

	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/repositories"
	"github.com/google/uuid"
)

type AttachmentService interface {
	GetByPublicID(pubId uuid.UUID) (*models.CardAttachment, error)
	Create(cardPublicId, userPublicID, filename string) (*models.CardAttachment, error)
	DeleteByPublicID(pubID uuid.UUID) error

	//gpt
	FindByCardID(cardPublicID string) ([]models.CardAttachment, error)
}

type attachmentService struct {
	AttachmentRepo repositories.AttachmentRepository
	cardRepo       repositories.CardRepository
	userRepo       repositories.UserRepository
}

func NewAttachmentService (
	attachmentRepo repositories.AttachmentRepository,
	cardRepo repositories.CardRepository,
	userRepo repositories.UserRepository,
) AttachmentService {
	return &attachmentService{attachmentRepo, cardRepo, userRepo}
}

func (s *attachmentService) GetByPublicID(pubId uuid.UUID) (*models.CardAttachment, error) {
	return s.AttachmentRepo.GetByPublicID(pubId)
}

func (s *attachmentService) Create(cardPublicId, userPublicID, filename string) (*models.CardAttachment, error) {
	card, err := s.cardRepo.FindByPublicID(cardPublicId)
	if err != nil {
		return nil, errors.New("card not found")
	}
	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	attach := &models.CardAttachment{
		PublicID: uuid.New(),
		CardID:   card.InternalID,
		UserID: user.InternalID,
		File: filename,
		CreatedAt: time.Now(),
	}

	if err := s.AttachmentRepo.Create(attach); err != nil {
		return nil, err
	}
	return attach, nil
}

func (s *attachmentService) DeleteByPublicID(pubID uuid.UUID) error {
	return s.AttachmentRepo.DeleteByPublicID(pubID)
}

func (s *attachmentService) FindByCardID(cardPublicID string) ([]models.CardAttachment, error) {
    return s.AttachmentRepo.FindByCardID(cardPublicID)
}