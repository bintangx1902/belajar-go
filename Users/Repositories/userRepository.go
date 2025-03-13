package Repositories

import (
	"TestApp/Models"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) FindUserById(userID uint) bool {
	var user Models.Users
	result := repo.DB.First(&user, userID)
	return !(result.Error == nil)
}

func (repo *UserRepository) FindUserByUsername(username string) bool {
	var user Models.Users
	result := repo.DB.Where("username = ?", username).First(&user)
	return result.Error == nil
}

func (repo *UserRepository) CreateUser(user *Models.Users) error {
	return repo.DB.Create(user).Error
}

func (repo *UserRepository) FindByUsernameAndPassword(username string, password string) (*Models.Users, error) {
	var user Models.Users
	err := repo.DB.Where(fmt.Sprintf("username = '%s' AND password = '%s'", username, password)).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
