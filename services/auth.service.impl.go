package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tonybobo/auth-template/models"
	"github.com/tonybobo/auth-template/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	AuthRepository models.AuthRepository
	ctx        context.Context
}

func NewAuthService(AuthRepository models.AuthRepository, ctx context.Context) AuthService {
	return &AuthServiceImpl{AuthRepository, ctx}
}

func (uc *AuthServiceImpl) SignInUser(user *models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = false
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	res, err := uc.collection.InsertOne(uc.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := uc.collection.Indexes().CreateOne(uc.ctx, index); err != nil {
		return nil, errors.New("cannot create index for email")
	}

	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}

	if err := uc.collection.FindOne(uc.ctx, query).Decode(&newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}
