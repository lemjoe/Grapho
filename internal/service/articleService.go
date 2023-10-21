package service

import (
	"fmt"
	"time"

	"github.com/lemjoe/md-blog/internal/models"
	"github.com/lemjoe/md-blog/internal/repository"
)

type articleService struct {
	repository  *repository.Repository
	fileService FileService
}

func NewArticleService(repository *repository.Repository, fileService FileService) *articleService {
	return &articleService{
		repository:  repository,
		fileService: fileService,
	}
}

// implement func ArticleService interface
func (a *articleService) CreateNewArticle(fileName, title string, author string, body []byte) (models.Article, error) {
	authorInfo, err := a.repository.User.GetUser(author)
	if err != nil {
		return models.Article{}, err
	}
	art, err := a.repository.Article.CreateArticle(models.Article{
		FileName:         fileName,
		Title:            title,
		Author:           authorInfo.FullName,
		AuthorId:         author,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		IsLocked:         false,
	})
	if err != nil {
		return models.Article{}, err
	}
	//write to file
	err = a.fileService.CreateNewFile(art.FileName, body)
	if err != nil {
		lockErr := a.repository.Article.LockArticleByFileName(art.FileName)
		if lockErr != nil {
			return models.Article{}, fmt.Errorf("unable to create file and lock article[%s]: \nerr: %w\nlockErr: %w", art.FileName, err, lockErr)
		}
		return models.Article{}, fmt.Errorf("unable to create article file[%s]: %w", art.FileName, lockErr)
	}
	return art, nil
}

func (a *articleService) DeleteArticle(fileName string) error {
	err := a.fileService.DeleteFile("articles/" + fileName)
	if err != nil {
		return fmt.Errorf("unable to delete file[%s]: %w", fileName, err)
	}
	err = a.repository.Article.DeleteArticleByFileName(fileName)
	if err != nil {
		return fmt.Errorf("unable to delete article[%s]: %w", fileName, err)
	}
	return nil
}

func (a *articleService) GetArticleInfo(fileName string) (models.Article, error) {
	art, err := a.repository.Article.GetArticleByFileName(fileName)
	if err != nil {
		return models.Article{}, err
	}
	return art, nil
}
func (a *articleService) GetArticleBody(fileName string) ([]byte, error) {
	file, err := a.fileService.ReadFile("articles/" + fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}
