package repositories

import (
	"github.com/adamabiyuu/project-management/config"
	"github.com/adamabiyuu/project-management/models"
)

type LabelRepository interface {
	FindByPublicID(publicID string) (*models.Label, error)
}

type labelRepository struct{}

func NewLabelRepository() LabelRepository {
	return &labelRepository{}
}

func (r *labelRepository) FindByPublicID(publicID string) (*models.Label, error) {
	var label models.Label

	if err := config.DB.
		Where("public_id = ?", publicID).
		First(&label).Error; err != nil {
		return nil, err
	}

	return &label, nil
}