package user

import (
	"fmt"
	"time"

	"github.com/lemjoe/Grapho/internal/models"
	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	q "github.com/ostafen/clover/v2/query"
)

type User struct {
	collectionName string
	db             *c.DB
}
type userSchema struct {
	UserName     string            `json:"user_name"`
	FullName     string            `json:"full_name"`
	Password     string            `json:"passwd"`
	Email        string            `json:"email"`
	IsAdmin      bool              `json:"is_admin"`
	IsWriter     bool              `json:"is_writer"`
	Id           string            `json:"_id"`
	LastLogin    time.Time         `json:"last_login"`
	CreationDate time.Time         `json:"creation_date"`
	Settings     map[string]string `json:"settings"`
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
	_, err := u.GetUserByUsername(user.UserName)
	if err == nil {
		return models.User{}, fmt.Errorf("user already exists")
	}
	doc := d.NewDocument()
	doc.Set("user_name", user.UserName)
	doc.Set("full_name", user.FullName)
	doc.Set("passwd", user.Password)
	doc.Set("email", user.Email)
	doc.Set("is_admin", user.IsAdmin)
	doc.Set("is_writer", user.IsWriter)
	doc.Set("creation_date", time.Now())
	doc.Set("last_login", time.Now())
	doc.Set("settings", user.Settings)
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
		IsWriter:     user.IsWriter,
		Id:           docId,
		LastLogin:    time.Now(),
		CreationDate: time.Now(),
		Settings:     user.Settings,
	}, nil
}

func (u *User) GetAllUsers() ([]models.User, error) {
	var users []models.User
	docs, err := u.db.FindAll(q.NewQuery(u.collectionName))
	if err != nil {
		return nil, fmt.Errorf("unable to find documents[%s]: %w", u.collectionName, err)
	}
	for _, doc := range docs {
		var user userSchema
		err := doc.Unmarshal(&user)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal document[%s]: %w", u.collectionName, err)
		}

		users = append(users, models.User{
			UserName:     user.UserName,
			FullName:     user.FullName,
			Password:     user.Password,
			Email:        user.Email,
			IsAdmin:      user.IsAdmin,
			IsWriter:     user.IsWriter,
			Id:           user.Id,
			LastLogin:    user.LastLogin,
			CreationDate: user.CreationDate,
			Settings:     user.Settings,
		})
	}
	return users, nil
}

// GetUser(username string) (models.User, error)
func (u *User) GetUserByUsername(username string) (models.User, error) {
	userRow, err := u.db.FindFirst(q.NewQuery(u.collectionName).Where(q.Field("user_name").Eq(username)))
	if err != nil {
		return models.User{}, err
	}
	if userRow == nil {
		return models.User{}, fmt.Errorf("user not found")
	}
	//fmt.Printf("userRow[%s]: %v\n", username, userRow)
	var user userSchema
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
		IsWriter:     user.IsWriter,
		Id:           user.Id,
		LastLogin:    user.LastLogin,
		CreationDate: user.CreationDate,
		Settings:     user.Settings,
	}, nil
}

// GetUserById(id string) (models.User, error)
func (u *User) GetUserById(id string) (models.User, error) {
	userRow, err := u.db.FindFirst(q.NewQuery(u.collectionName).Where(q.Field("_id").Eq(id)))
	if err != nil {
		return models.User{}, err
	}
	if userRow == nil {
		return models.User{}, fmt.Errorf("user not found")
	}
	//fmt.Printf("userRow[%s]: %v\n", username, userRow)
	var user userSchema
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
		IsWriter:     user.IsWriter,
		Id:           user.Id,
		LastLogin:    user.LastLogin,
		CreationDate: user.CreationDate,
		Settings:     user.Settings,
	}, nil

}

func (u *User) ChangeUserPassword(id string, passwd string) error {
	err := u.db.UpdateById(u.collectionName, id, func(doc *d.Document) *d.Document {
		doc.Set("passwd", passwd)
		return doc
	})
	return err
}

// ChangeUserSettings(id string, settings map[string]string) error
func (u *User) ChangeUserSettings(id string, settings map[string]string) error {
	err := u.db.UpdateById(u.collectionName, id, func(doc *d.Document) *d.Document {
		doc.Set("settings", settings)
		return doc
	})
	return err
}

func (u *User) UpdateUserData(id string, fullname string, email string, isadmin bool, iswriter bool) error {
	err := u.db.UpdateById(u.collectionName, id, func(doc *d.Document) *d.Document {
		doc.Set("full_name", fullname)
		doc.Set("email", email)
		doc.Set("is_admin", isadmin)
		doc.Set("is_writer", iswriter)
		return doc
	})
	return err
}
