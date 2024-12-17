package users

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/users/model"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (model.User, error)
	GetManyByIDS(ctx context.Context, input GetManyInput) ([]model.User, error)
	Update(ctx context.Context, id string, input UpdateInput) (model.User, error)
}

type GetManyInput struct {
	Page       int
	PerPage    int
	IDs        []string
	IsBotAdmin *bool
	IsBanned   *bool
}

type UpdateInput struct {
	IsBanned          *bool
	IsBotAdmin        *bool
	ApiKey            *string
	HideOnLandingPage *bool
	TokenID           *string
}