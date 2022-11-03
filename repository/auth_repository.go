package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tonybobo/auth-template/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *authCollection) UpdateUserById(ctx context.Context, id primitive.ObjectID, field, value string) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := r.DB.UpdateOne(ctx, query, update)
	return result, err
}

func (r *authCollection) UpdateOne(ctx context.Context, field string, value interface{}) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: field, Value: value}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}

	result, err := r.DB.UpdateOne(ctx, query, update)

	return result, err
}

func (r *authCollection) ResetPasswordToken(ctx context.Context, email, passwordResetToken string) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: "email", Value: strings.ToLower(email)}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "passwordResetToken", Value: passwordResetToken},
			{Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)},
		}},
	}
	result, err := r.DB.UpdateOne(ctx, query, update)

	return result, err
}

func (r *authCollection) VerifyEmail(ctx context.Context, verificationCode string) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}},
		{Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}

	result, err := r.DB.UpdateOne(ctx, query, update)

	return result, err
}

func (r *authCollection) ClearResetPasswordToken(ctx context.Context, token, password string) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: "passwordResetToken", Value: token}, {Key: "passwordResetAt", Value: bson.D{{Key: "$gt", Value: time.Now()}}}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "password", Value: password}}},
		{Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}

	result, err := r.DB.UpdateOne(ctx, query, update)

	return result, err
}

func (r *authCollection) SignUpUser(ctx context.Context, user *models.SignUpInput) (*models.DBResponse, error) {
	result, err := r.DB.InsertOne(ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := r.DB.Indexes().CreateOne(ctx, index); err != nil {
		return nil, errors.New("cannot create index for email")
	}

	var newUser *models.DBResponse
	query := bson.M{"_id": result.InsertedID}

	if err := r.DB.FindOne(ctx, query).Decode(&newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}
