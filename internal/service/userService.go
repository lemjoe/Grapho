package service

import (
	"time"

	"github.com/lemjoe/md-blog/internal/models"
	"github.com/lemjoe/md-blog/internal/repository/repotypes"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repository *repotypes.Repository
}

func NewUserService(repository *repotypes.Repository) *userService {
	return &userService{
		repository: repository,
	}
}

// implement func ArticleService interface
func (u *userService) CreateNewUser(username string, fullname string, password string, email string, isadmin bool) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return models.User{}, err
	}

	user, err := u.repository.User.CreateUser(models.User{
		// FileName:         fileName,
		UserName:     username,
		FullName:     fullname,
		Password:     string(hash),
		Email:        email,
		IsAdmin:      isadmin,
		CreationDate: time.Now(),
		LastLogin:    time.Unix(int64(0), int64(0)),
		Settings:     DefaultUserSettings,
	})
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userService) GetUserById(id string) (models.User, error) {
	user, err := u.repository.User.GetUserById(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userService) GetUserByName(username string) (models.User, error) {
	user, err := u.repository.User.GetUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userService) ChangeUserPassword(id string, passwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), 10)
	if err != nil {
		return err
	}
	err = u.repository.User.ChangeUserPassword(id, string(hash))
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) ChangeUserSettings(id string, settings map[string]string) error {
	err := u.repository.User.ChangeUserSettings(id, settings)
	if err != nil {
		return err
	}
	return nil
}
