package repotypes

import "github.com/lemjoe/Grapho/internal/models"

type User interface {
	CreateUser(user models.User) (models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserById(id string) (models.User, error)
	ChangeUserPassword(string, string) error
	ChangeUserSettings(id string, settings map[string]string) error
	UpdateUserData(id string, fullname string, email string, isadmin bool, iswriter bool) error
}
type Article interface {
	CreateArticle(article models.Article) (models.Article, error)
	GetAllArticles() ([]models.Article, error) //todo add pagination
	GetArticleById(id string) (models.Article, error)
	DeleteArticleById(id string) error
	UpdateArticleById(id string) error
	LockArticleById(id string) error
}
type Repository struct {
	User    User
	Article Article
}
