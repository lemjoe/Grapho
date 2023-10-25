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
	usScheme := UserScheme{
		UserName:     user.UserName,
		FullName:     user.FullName,
		Password:     user.Password,
		Email:        user.Email,
		IsAdmin:      user.IsAdmin,
		LastLogin:    time.Now(),
		CreationDate: time.Now(),
	}
	res, err := u.ct.InsertOne(context.TODO(), bson.M{
		"user_name":     usScheme.UserName,
		"full_name":     usScheme.FullName,
		"passwd":        usScheme.Password,
		"email":         usScheme.Email,
		"is_admin":      usScheme.IsAdmin,
		"last_login":    usScheme.LastLogin,
		"creation_date": usScheme.CreationDate,
	})
	if err != nil {
		return models.User{}, err
	}
	usScheme.Id = res.InsertedID.(primitive.ObjectID)
	return models.User{
		UserName:     usScheme.UserName,
		FullName:     usScheme.FullName,
		Password:     usScheme.Password,
		Email:        usScheme.Email,
		IsAdmin:      usScheme.IsAdmin,
		Id:           usScheme.Id.Hex(),
		LastLogin:    usScheme.LastLogin,
		CreationDate: usScheme.CreationDate,
	}, nil
}

// GetUser(username string) (models.User, error)
func (u *User) GetUserByUsername(username string) (models.User, error) {
	var findedUser UserScheme
	err := u.ct.FindOne(context.TODO(), bson.M{"user_name": username}).Decode(&findedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, fmt.Errorf("user not found")
		}

		return models.User{}, err
	}
	//fmt.Printf("findedUser: %+v\n", findedUser)
	return models.User{
		UserName:     findedUser.UserName,
		FullName:     findedUser.FullName,
		Password:     findedUser.Password,
		Email:        findedUser.Email,
		IsAdmin:      findedUser.IsAdmin,
		Id:           findedUser.Id.Hex(),
		LastLogin:    findedUser.LastLogin,
		CreationDate: findedUser.CreationDate,
	}, nil
}

// GetUserById(id string) (models.User, error)
func (u *User) GetUserById(id string) (models.User, error) {
	var findedUser UserScheme
	err := u.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, err
	}
	return models.User{
		UserName:     findedUser.UserName,
		FullName:     findedUser.FullName,
		Password:     findedUser.Password,
		Email:        findedUser.Email,
		IsAdmin:      findedUser.IsAdmin,
		Id:           findedUser.Id.Hex(),
		LastLogin:    findedUser.LastLogin,
		CreationDate: findedUser.CreationDate,
	}, nil
}
