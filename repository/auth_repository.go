package repository

import (
	"context"
	"strings"

	"github.com/tonybobo/auth-template/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type authCollection struct {
	DB *mongo.Collection
}

func NewAuthRepository(db *mongo.Collection) models.AuthRepository {
	return &authCollection{DB: db}
}

func (r *authCollection) FindUserById(ctx context.Context, id primitive.ObjectID) (*models.DBResponse, error) {
	var user *models.DBResponse
	query := bson.M{"_id": id}
	if err := r.DB.FindOne(ctx, query).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil
}

func (r *authCollection) FindUserByEmail(ctx context.Context, email string) (*models.DBResponse, error) {
	var user *models.DBResponse
	query := bson.M{"email": strings.ToLower(email)}
	if err := r.DB.FindOne(ctx, query).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}
	return user, nil
}
