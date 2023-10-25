package article

import (
	"context"
	"fmt"
	"time"

	"github.com/lemjoe/md-blog/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type articleSchema struct {
	Title            string             `json:"article_title"`
	Author           string             `json:"article_author"`
	AuthorId         string             `json:"author_id"`
	CreationDate     time.Time          `json:"creation_date"`
	ModificationDate time.Time          `json:"modification_date"`
	IsLocked         bool               `json:"is_locked"`
	Id               primitive.ObjectID `json:"_id"`
}
type Article struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) (*Article, error) {
	collectionName := "articles"
	names, err := driver.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return &Article{}, fmt.Errorf("unable to list collections: %w", err)
	}
	//	fmt.Printf("names: %+v\n", names)
	if !slices.Contains(names, collectionName) {
		command := bson.M{"create": collectionName}
		var result bson.M
		if err := driver.RunCommand(context.TODO(), command).Decode(&result); err != nil {
			return &Article{}, fmt.Errorf("unable to create collection[%s]: %w", collectionName, err)
		}
	}

	return &Article{
		ct: driver.Collection(collectionName),
	}, nil
}

// add interface empty methods
// CreateArticle(article models.Article) (models.Article, error)
//
//	GetAllArticles() ([]models.Article, error) //todo add pagination
//	GetArticleById(id string) (models.Article, error)
//	DeleteArticleById(id string) error
//	UpdateArticleById(id string) error
//	LockArticleById(id string) error
func (a *Article) CreateArticle(article models.Article) (models.Article, error) {
	return models.Article{}, nil
}
func (a *Article) GetAllArticles() ([]models.Article, error) {
	return []models.Article{}, nil
}
func (a *Article) GetArticleById(id string) (models.Article, error) {
	return models.Article{}, nil
}
func (a *Article) DeleteArticleById(id string) error {
	return nil
}
func (a *Article) UpdateArticleById(id string) error {
	return nil
}
func (a *Article) LockArticleById(id string) error {
	return nil
}
