package repository

import (
	"github.com/lemjoe/md-blog/internal/models"
	"github.com/lemjoe/md-blog/internal/repository/cloverdb/article"
	"github.com/lemjoe/md-blog/internal/repository/cloverdb/user"
	"github.com/ostafen/clover/v2"
)

type User interface {
	CreateUser(user models.User) (models.User, error)
	GetUser(username string) (models.User, error)
}
type Article interface {
	CreateArticle(article models.Article) (models.Article, error)
	GetAllArticles() ([]models.Article, error) //todo add pagination
	GetArticleByFileName(fileName string) (models.Article, error)
	DeleteArticleByFileName(fileName string) error
	UpdateArticleByFileName(fileName string) error
	LockArticleByFileName(fileName string) error
}
type Repository struct {
	User    User
	Article Article
}

func NewRepository(db *clover.DB) (*Repository, error) {
	user, err := user.Init(db)
	if err != nil {
		return nil, err
	}
	article, err := article.Init(db)
	if err != nil {
		return nil, err
	}
	return &Repository{
		User:    user,
		Article: article,
	}, nil
}
