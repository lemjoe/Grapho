package user

import (
	"fmt"
	"time"

	"github.com/lemjoe/md-blog/internal/models"
	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	q "github.com/ostafen/clover/v2/query"
)

type User struct {
	collectionName string
	db             *c.DB
}
type userSchema struct {
	UserName     string    `json:"user_name"`
	FullName     string    `json:"full_name"`
	Password     string    `json:"passwd"`
	Email        string    `json:"email"`
	IsAdmin      bool      `json:"is_admin"`
	Id           string    `json:"id"`
	LastLogin    time.Time `json:"last_login"`
	CreationDate time.Time `json:"creation_date"`
}

func Init(db *c.DB) (*User, error) {
	collection := User{
		collectionName: "users",
		db:             db,
	}
	//check if collection already exists
	exists, err := db.HasCollection(collection.collectionName)
	if err != nil {
		return nil, fmt.Errorf("unable to check if collection[%s] exists: %w", collection.collectionName, err)
	}
	if !exists {
		err := db.CreateCollection(collection.collectionName)
		if err != nil {
			return nil, fmt.Errorf("unable to create collection[%s]: %w", collection.collectionName, err)
		}
	}
	return &collection, nil
}

// CreateUser(user models.User) (models.User, error)
func (u *User) CreateUser(user models.User) (models.User, error) {
	//check if user already exists
	_, err := u.GetUser(user.UserName)
	if err == nil {
		return models.User{}, fmt.Errorf("user already exists")
	}
	doc := d.NewDocument()
	doc.Set("user_name", user.UserName)
	doc.Set("full_name", user.FullName)
	doc.Set("passwd", user.Password)
	doc.Set("email", user.Email)
	doc.Set("is_admin", user.IsAdmin)
	doc.Set("creation_date", time.Now())
	doc.Set("last_login", time.Now())
	docId, err := u.db.InsertOne(u.collectionName, doc)
	if err != nil {
		return models.User{}, fmt.Errorf("unable to insert document[%s]: %w", u.collectionName, err)
	}
	return models.User{
		UserName:     user.UserName,
		FullName:     user.FullName,
		Password:     user.Password,
		Email:        user.Email,
		IsAdmin:      user.IsAdmin,
		Id:           docId,
		LastLogin:    time.Now(),
		CreationDate: time.Now(),
	}, nil
}

// GetUser(username string) (models.User, error)
func (u *User) GetUser(username string) (models.User, error) {
	userRow, err := u.db.FindFirst(q.NewQuery(u.collectionName).Where(q.Field("user_name").Eq(username)))
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	err = userRow.Unmarshal(&user)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		UserName:     user.UserName,
		FullName:     user.FullName,
		Password:     user.Password,
		Email:        user.Email,
		IsAdmin:      user.IsAdmin,
		Id:           user.Id,
		LastLogin:    user.LastLogin,
		CreationDate: user.CreationDate,
	}, nil
}
