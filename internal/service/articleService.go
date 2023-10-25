package service

import (
	"fmt"
	"time"

	"github.com/lemjoe/md-blog/internal/models"
	"github.com/lemjoe/md-blog/internal/repository/repotypes"
)

type articleService struct {
	repository  *repotypes.Repository
	fileService FileService
}

func NewArticleService(repository *repotypes.Repository, fileService FileService) *articleService {
	return &articleService{
		repository:  repository,
		fileService: fileService,
	}
}

// implement func ArticleService interface
func (a *articleService) CreateNewArticle(title string, author string, body []byte) (models.Article, error) {
	// fileName := hash.GetHash(body)
	authorInfo, err := a.repository.User.GetUserByUsername(author)
	if err != nil {
		return models.Article{}, err
	}
	art, err := a.repository.Article.CreateArticle(models.Article{
		// FileName:         fileName,
		Title:            title,
		Author:           authorInfo.FullName,
		AuthorId:         authorInfo.Id,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		IsLocked:         false,
	})
	if err != nil {
		return models.Article{}, err
	}
	//write to file
	err = a.fileService.CreateNewFile("articles/"+art.Id, body)
	if err != nil {
		lockErr := a.repository.Article.LockArticleById(art.Id)
		if lockErr != nil {
			return models.Article{}, fmt.Errorf("unable to create file and lock article[%s]: \nerr: %w\nlockErr: %w", art.Id, err, lockErr)
		}
		return models.Article{}, fmt.Errorf("unable to create article file[%s]: %w", art.Id, lockErr)
	}
	return art, nil
}

func (a *articleService) DeleteArticle(id string) error {
	err := a.fileService.DeleteFile("articles/" + id)
	if err != nil {
		return fmt.Errorf("unable to delete file[%s]: %w", id, err)
	}
	err = a.repository.Article.DeleteArticleById(id)
	if err != nil {
		return fmt.Errorf("unable to delete article[%s]: %w", id, err)
	}
	return nil
}

func (a *articleService) GetArticleInfo(id string) (models.Article, error) {
	art, err := a.repository.Article.GetArticleById(id)
	if err != nil {
		return models.Article{}, err
	}
	return art, nil
}
func (a *articleService) GetArticleBody(id string) ([]byte, error) {
	file, err := a.fileService.ReadFile("articles/" + id)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetArticlesList() ([]models.Article, error)
func (a *articleService) GetArticlesList() ([]models.Article, error) {
	arts, err := a.repository.Article.GetAllArticles()
	if err != nil {
		return nil, err
	}
	return arts, nil
}

// UpdateArticle(fileName string) error
func (a *articleService) UpdateArticle(fileName string) error {
	return a.repository.Article.UpdateArticleById(fileName)
}
