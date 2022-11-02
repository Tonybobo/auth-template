package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRepository interface {
	FindUserById(ctx context.Context, id primitive.ObjectID) (*DBResponse, error)
	FindUserByEmail(ctx context.Context, email string) (*DBResponse, error)
}
