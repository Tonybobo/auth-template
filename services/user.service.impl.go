package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tonybobo/auth-template/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid User ID")
	}

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	if err := us.collection.FindOne(us.ctx, query).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil

}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse
	query := bson.M{"email": strings.ToLower(email)}

	if err := us.collection.FindOne(us.ctx, query).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) UpdateUserById(id, field, value string) (*models.DBResponse, error) {
	userId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := us.collection.UpdateOne(us.ctx, query, update)

	fmt.Print(result.ModifiedCount)

	if err != nil {
		fmt.Print(err)
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}

func (us *UserServiceImpl) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {
	query := bson.D{{Key: field, Value: value}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := us.collection.UpdateOne(us.ctx, query, update)

	fmt.Print(result.ModifiedCount)
	if err != nil {
		fmt.Print(err)
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}
