package graph

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"graphql/Models"
	"graphql/graph/model"
)

type Resolver struct {
	DB *gorm.DB
}

func (r *Resolver) User(ctx context.Context) ([]*model.User, error) {
	var User []*model.User
	err := r.DB.Find(&User).Error
	if err != nil {
		return nil, err
	}
	return User, nil
}
func (r *Resolver) CreateUser(ctx context.Context, name string, email string) (*Models.User, error) {
	user := &Models.User{
		ID:    0,
		Name:  name,
		Email: email,
	}

	if err := r.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
