package service

import "github.com/lemjoe/Grapho/internal/models"

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
	ListFolder(path string) ([][]byte, error)
}
type MigrationService interface {
	Migrate(string) error
}
type UserService interface {
	CreateNewUser(string, string, string, string, bool) (models.User, error)
	GetUserById(string) (models.User, error)
	GetUserByName(string) (models.User, error)
	ChangeUserPassword(string, string) error
	ChangeUserSettings(string, map[string]string) error
	GetUsersList() ([]models.User, error)
	UpdateUserData(string, string, string, bool) error
}
