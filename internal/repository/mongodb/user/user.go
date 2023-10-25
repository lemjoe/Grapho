package user

import (
	"context"
	"time"

	"fmt"

	"github.com/lemjoe/md-blog/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type UserScheme struct {
	UserName     string             `json:"user_name"`
	FullName     string             `json:"full_name"`
	Password     string             `json:"passwd"`
	Email        string             `json:"email"`
	IsAdmin      bool               `json:"is_admin"`
	Id           primitive.ObjectID `json:"_id"`
	LastLogin    time.Time          `json:"last_login"`
	CreationDate time.Time          `json:"creation_date"`
}
type User struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) (*User, error) {
	collectionName := "users"
	names, err := driver.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return &User{}, fmt.Errorf("unable to list collections: %w", err)
	}
	//	fmt.Printf("names: %+v\n", names)
	if !slices.Contains(names, collectionName) {
		command := bson.M{"create": collectionName}
		var result bson.M
		if err := driver.RunCommand(context.TODO(), command).Decode(&result); err != nil {
			return &User{}, fmt.Errorf("unable to create collection[%s]: %w", collectionName, err)
		}
	}

	return &User{
		ct: driver.Collection(collectionName),
	}, nil
}
func (u *User) CreateUser(user models.User) (models.User, error) {
	return models.User{}, nil
}

// GetUser(username string) (models.User, error)
func (u *User) GetUser(username string) (models.User, error) {
	return models.User{}, nil
}
