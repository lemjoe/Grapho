package user

import (
	"context"
	"time"

	"fmt"

	"github.com/lemjoe/Grapho/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type UserScheme struct {
	UserName     string             `bson:"user_name"`
	FullName     string             `bson:"full_name"`
	Password     string             `bson:"passwd"`
	Email        string             `bson:"email"`
	IsAdmin      bool               `bson:"is_admin"`
	Id           primitive.ObjectID `bson:"_id"`
	LastLogin    time.Time          `bson:"last_login"`
	CreationDate time.Time          `bson:"creation_date"`
	Settings     map[string]string  `bson:"settings"`
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
		Settings:     user.Settings,
	}
	res, err := u.ct.InsertOne(context.TODO(), bson.M{
		"user_name":     usScheme.UserName,
		"full_name":     usScheme.FullName,
		"passwd":        usScheme.Password,
		"email":         usScheme.Email,
		"is_admin":      usScheme.IsAdmin,
		"last_login":    usScheme.LastLogin,
		"creation_date": usScheme.CreationDate,
		"settings":      usScheme.Settings,
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
		Settings:     user.Settings,
	}, nil
}

func (u *User) GetAllUsers() ([]models.User, error) {
	var findedUsers []models.User
	cur, err := u.ct.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var usr UserScheme
		err := cur.Decode(&usr)
		if err != nil {
			return nil, err
		}
		findedUsers = append(findedUsers, models.User{
			UserName:     usr.UserName,
			FullName:     usr.FullName,
			Password:     usr.Password,
			Email:        usr.Email,
			IsAdmin:      usr.IsAdmin,
			Id:           usr.Id.Hex(),
			LastLogin:    usr.LastLogin,
			CreationDate: usr.CreationDate,
			Settings:     usr.Settings,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedUsers, nil
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
		Settings:     findedUser.Settings,
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
		Settings:     findedUser.Settings,
	}, nil
}

func (u *User) ChangeUserPassword(id string, passwd string) error {
	usrObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = u.ct.UpdateOne(context.TODO(), bson.M{"_id": usrObjId}, bson.M{"$set": bson.M{
		"passwd": passwd,
	}})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ChangeUserSettings(id string, settings map[string]string) error {
	usrObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = u.ct.UpdateOne(context.TODO(), bson.M{"_id": usrObjId}, bson.M{"$set": bson.M{
		"settings": settings,
	}})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateUserData(id string, fullname string, email string, isadmin bool) error {
	usrObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = u.ct.UpdateOne(context.TODO(), bson.M{"_id": usrObjId}, bson.M{"$set": bson.M{
		"full_name": fullname,
		"email":     email,
		"is_admin":  isadmin,
	}})
	if err != nil {
		return err
	}
	return nil
}
