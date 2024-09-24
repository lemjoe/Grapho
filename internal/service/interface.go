package service

import "github.com/lemjoe/md-blog/internal/models"

type ArticleService interface {
	CreateNewArticle(string, string, []byte) (models.Article, error)
	DeleteArticle(fileName string) error
	UpdateArticle(fileName string) error
	//GetArticle(fileName string) (string, error)
	GetArticleInfo(fileName string) (models.Article, error)
	GetArticleBody(fileName string) ([]byte, error)
	GetArticlesList() ([]models.Article, error)
}
type FileService interface {
	ReadFile(fileName string) ([]byte, error)
	CreateNewFile(path string, body []byte) error
	WriteFile(path string, body []byte) error
	CreateFolder(path string) error
	DeleteFile(path string) error
}
type MigrationService interface {
	Migrate(string) error
}
type UserService interface {
	CreateNewUser(string, string, string, string, bool) (models.User, error)
	GetUserById(string) (models.User, error)
	GetUserByName(string) (models.User, error)
	ChangeUserSettings(string, map[string]string) error
}
