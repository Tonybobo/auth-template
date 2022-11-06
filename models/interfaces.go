package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	FindUserById(ctx context.Context, id primitive.ObjectID) (*DBResponse, error)
	FindUserByEmail(ctx context.Context, email string) (*DBResponse, error)
	UpdateOne(ctx context.Context, field string, value interface{}) (*mongo.UpdateResult, error)
	ResetPasswordToken(ctx context.Context, email, passwordResetToken string) (*mongo.UpdateResult, error)
	VerifyEmail(ctx context.Context, verificationCode string) (*mongo.UpdateResult, error)
	ClearResetPasswordToken(ctx context.Context, token, password string) (*mongo.UpdateResult, error)
	SignUpUser(ctx context.Context, user *SignUpInput) (*DBResponse, error, string)
}
