package storage

import (
	"context"
	"errors"
	"talk/internal/domain"
)

type Storage interface {
	GetUserByExternalUserID(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

var ErrNotFound = errors.New("not found")
