package repositories

import (
	"fmt"
	"path/filepath"

	"github.com/adamabiyuu/project-management/config"
	"github.com/adamabiyuu/project-management/models"
	"gorm.io/gorm"
)

type CardRepository interface {
	Create(card *models.Card) error
	Update(card *models.Card) error
	Delete(id uint) error
	FindByID(id uint) (*models.Card, error)
	FindByPublicID(publicID string) (*models.Card, error)
	FindByListID (listID string)([]*models.Card, error)
}

type cardRepository struct {

}

func NewCardRepository() CardRepository {
	return &cardRepository{}
}

func (r *cardRepository) Create(card *models.Card) error {
	return config.DB.Create(card).Error
}

func (r *cardRepository) Update(card *models.Card) error {
	return config.DB.Save(card).Error
}

func (r *cardRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Card{}, id).Error
}

func (r *cardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	err := config.DB.Preload("Labels").Preload("Assigness").First(&card,id).Error

	return &card, err
}

func (r *cardRepository) FindByPublicID(publicID string) (*models.Card, error) {
	var card models.Card

	// Mengambil relasi Assignees beserta data User dari setiap assignee
	// Misalnya Card -> Assignees -> User
	if err := config.DB.Preload("Assignees.User", func (tx *gorm.DB) *gorm.DB {
		// Hanya mengambil kolom tertentu dari tabel users
		// agar response lebih ringan dan tidak mengirim data yang tidak diperlukan
		return tx.Select("internal_id", "public_id", "name", "email")
	}).Preload("Attachments").Where("public_id = ?", publicID).First(&card).Error; err != nil {
		return nil, err
	}

	// Mengambil base URL aplikasi dari file .env
	// Contoh:
	// http://localhost:3030
	baseUrl := config.AppConfig.APPURL

	// Melakukan perulangan untuk setiap attachment yang dimiliki card
	for i := range card.Attachments {
		// Membuat URL yang bisa diakses oleh browser
		//
		// Misalnya:
		// File di database:
		// uploads/card/jwt.pdf
		//
		// Menjadi:
		// http://localhost:3030/files/jwt.pdf
		card.Attachments[i].FileURL = fmt.Sprintf("%s/files/%s",
		baseUrl,
		filepath.Base(card.Attachments[i].File),
	)
	}

	return &card, nil
}