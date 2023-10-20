package service

type ArticleService interface {
	CreateNewArticle(fileName, title string) error
	DeleteArticle(fileName string) error
	//	UpdateArticle(fileName string) error
	GetArticle(fileName string) (string, error)
}
