package repotypes

import "github.com/lemjoe/md-blog/internal/models"

type User interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserById(id string) (models.User, error)
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
