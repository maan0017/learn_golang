package repos

import (
	"context"
	"go/mongo-db/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
	}
}

func (r *UserRepo) GetUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []*models.User

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}) // -1 = descending, 1 = ascending
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepo) GetUserById(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	cursor := r.collection.FindOne(ctx, filter)

	var user models.User

	if err := cursor.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) CreateUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
