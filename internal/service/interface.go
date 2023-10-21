package service

type ArticleService interface {
	CreateNewArticle(fileName, title string, author string, body []byte) error
	DeleteArticle(fileName string) error
	//	UpdateArticle(fileName string) error
	GetArticle(fileName string) (string, error)
}
type FileService interface {
	ReadFile(fileName string) ([]byte, error)
	CreateNewFile(path string, body []byte) error
	WriteFile(path string, body []byte) error
	CreateFolder(path string) error
	DeleteFile(path string) error
}
