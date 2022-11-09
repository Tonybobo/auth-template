package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpInput struct {
	Name            string    `json:"name" bson:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	Role            string    `json:"role" bson:"role"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type DBResponse struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Name            string             `json:"name" bson:"name" binding:"required"`
	Email           string             `json:"email" bson:"email" binding:"required"`
	Password        string             `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string             `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	Role            string             `json:"role" bson:"role"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type ForgetPasswordInput struct {
	Email string `json:"email" bson:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
}

type AuthServiceResponse struct {
	User               *DBResponse
	Status             string
	Err                error
	Message            string
	StatusCode         int
	AccessToken        string
	RefreshAccessToken string
}

func FilteredResponse(result *DBResponse) UserResponse {
	return UserResponse{
		ID:        result.ID,
		Name:      result.Name,
		Email:     result.Email,
		Role:      result.Role,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}
}
