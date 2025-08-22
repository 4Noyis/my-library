package repositories

import (
	"context"
	"time"

	"github.com/4Noyis/my-library/internal/database"
	"github.com/4Noyis/my-library/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository struct {
	collection string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: "users",
	}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true

	_, err := database.Collection(ur.collection).InsertOne(ctx, user)
	return err
}

func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	filter := bson.D{{"username", username}}

	//err := UserCollection(&mongo.Client{}).FindOne(ctx, filter).Decode(&user)
	err := database.Collection(ur.collection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	filter := bson.D{{"email", email}}

	err := database.Collection(ur.collection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	filter := bson.D{{"_id", id}}

	err := database.Collection(ur.collection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(id primitive.ObjectID, updates bson.D) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"_id", id}}
	updates = append(updates, bson.E{"updated_at", time.Now()})
	update := bson.D{{"$set", updates}}

	_, err := database.Collection(ur.collection).UpdateOne(ctx, filter, update)
	return err
}
