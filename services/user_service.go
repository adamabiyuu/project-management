package services

import (
	"errors"

	"github.com/adamabiyuu/project-management/models"
	"github.com/adamabiyuu/project-management/repositories"
	"github.com/adamabiyuu/project-management/utils"
	"github.com/google/uuid"
)
// kontrak
type UserService interface {
	Register(user *models.User) error
}
// cetakan atau design blueprint
type userService struct {
	repo repositories.UserRepository
}
// cara buat object
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

// Dalam Kasusmu

// Kalau ditulis seperti ini:

// type UserService interface {...}

// type userService struct {...}

// func NewUserService(...) UserService {
// 	return &userService{repo}
// }

// func (s *userService) Register(...) error {
// 	...
// }

// bacanya menjadi:

// 1. Ini kontraknya (interface)
// 2. Ini bentuk object-nya (struct)
// 3. Ini cara membuat object-nya (constructor)
// 4. Ini kemampuan object-nya (method)


// Langkah 1: Interface
// type Animal interface {
// 	Sound()
// }

// Artinya:

// Siapa pun yang punya method Sound() bisa dianggap sebagai Animal.

// Langkah 2: Struct
// type Cat struct{}
// Langkah 3: Method
// func (c *Cat) Sound() {
// 	fmt.Println("Meow")
// }

// Sekarang *Cat punya method:

// Sound()
// Pertanyaan

// Apakah *Cat memenuhi interface:

// type Animal interface {
// 	Sound()
// }

// Jawabannya:

// ✅ Ya.

// Karena interface meminta:

// Sound()

// dan *Cat memiliki:

// Sound()

