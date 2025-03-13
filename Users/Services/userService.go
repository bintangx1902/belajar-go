package Services

import (
	"TestApp/Models"
	"TestApp/Users/Repositories"
	"TestApp/Utils"
	"errors"
	"fmt"
)

type UserService struct {
	Repo *Repositories.UserRepository
}

func NewUsersService(repo *Repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (service *UserService) CreateUser(user *Models.Users) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	err := service.Repo.FindUserByUsername(user.Username)
	if err {
		return fmt.Errorf("Error : %v", err)
	}
	return service.Repo.CreateUser(user)
}

func (service *UserService) Login(username, password string) (string, error) {
	user, err := service.Repo.FindByUsernameAndPassword(username, password)
	if err != nil {
		return "", errors.New("user not found")
	}

	if password != user.Password {
		return "", errors.New("invalid Credentials")
	}

	token, err := Utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
