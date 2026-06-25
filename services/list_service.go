package services

import (
	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/repositories"
	"github.com/google/uuid"
)

type listService struct {
	listRepo repositories.ListRepository
	boardRepo repositories.BoardRepository
	ListPosRepo repositories.ListPositionRepository
}

type ListWithOrder struct {
	Position []uuid.UUID
	Lists []models.List
}

type ListService interface {
	GetByBoardID(boardPublicID string) (*ListWithOrder, error)
	GetByID (id uint) (*models.List, error)
	GetByPublicID (publicID string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPublicID string, positions []uuid.UUID) error
}