package repositories

import (
	"github.com/adamabiyuu/project-management/config"
	"github.com/adamabiyuu/project-management/models"
	"github.com/google/uuid"
)

type ListRepository interface {
}

type listRepository struct {
}

func NewListRepository() ListRepository {
	return &listRepository{}
}

func (r *listRepository) Create(list *models.List) error {
	return config.DB.Create(list).Error
}

func (r *listRepository) Update(list *models.List) error {
	return config.DB.Model(&models.List{}).Where("public_id = ?", list.PublicId).
	Updates(map[string]interface{}{
		"tittle": list.Title,
	}).Error
}

func (r *listRepository) Delete(id uint) error {
	return config.DB.Delete(&models.List{}, id).Error
}

func (r *listRepository) UpdatePosition(boardPublicID string, position []string) error {
	return config.DB.Model(&models.ListPosition{}).
	Where("board_internal_id = (Select internal_id FROM boards WHERE public_id = ?)", boardPublicID).
	Update("list_oerder", position).Error
}

func (r *listRepository) GetCardPosition(listPublicID string) ([]uuid.UUID, error) {
	var position models.CardPosition
	err := config.DB.Joins("JOIN lists ON list.internal_id = card_positions.list_internal_id").
	Where("list.public_id = ?", listPublicID).Error
	return position.CardOrder, err
}