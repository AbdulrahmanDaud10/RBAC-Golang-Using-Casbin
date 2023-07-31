package repository

import (
	"log"

	"github.com/AbdulrahmanDaud10/RBAC-Casbin-Golang/pkg/api"
	"gorm.io/gorm"
)

// UserRepository : represent the user's repository contract
type UserRepository interface {
	AddUser(api.User) (api.User, error)
	GetUser(int) (api.User, error)
	GetByEmail(string) (api.User, error)
	GetAllUser() ([]api.User, error)
	UpdateUser(api.User) (api.User, error)
	DeleteUser(api.User) (api.User, error)
	Migrate() error
}

type userRepository struct {
	DB *gorm.DB
}

// NewUserRepository -> returns new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u userRepository) Migrate() error {
	log.Print("[UserRepository]...Migrate")
	return u.DB.AutoMigrate(&api.User{})
}

func (u userRepository) GetUser(id int) (user api.User, err error) {
	return user, u.DB.First(&user, id).Error
}

func (u userRepository) GetByEmail(email string) (user api.User, err error) {
	return user, u.DB.First(&user, "email=?", email).Error
}

func (u userRepository) GetAllUser() (users []api.User, err error) {
	return users, u.DB.Find(&users).Error
}

func (u userRepository) AddUser(user api.User) (api.User, error) {
	return user, u.DB.Create(&user).Error
}

func (u userRepository) UpdateUser(user api.User) (api.User, error) {
	if err := u.DB.First(&user, user.ID).Error; err != nil {
		return user, err
	}
	return user, u.DB.Model(&user).Updates(&user).Error
}

func (u userRepository) DeleteUser(user api.User) (api.User, error) {
	if err := u.DB.First(&user, user.ID).Error; err != nil {
		return user, err
	}
	return user, u.DB.Delete(&user).Error
}
