package services

import (
	"errors"

	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/repositories"
	"github.com/adamabiyuu/project-management/utils"
	"github.com/google/uuid"
)

type UserService interface {
	Register(user *models.User) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user *models.User) error {
	// kita harus mengecek email yang terdaftar atau belum
	// hash password
	// set role
	// simpan user

	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}
	hased, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hased
	user.Role = "user"
	user.PublicID = uuid.New()
	return s.repo.Create(user)
}