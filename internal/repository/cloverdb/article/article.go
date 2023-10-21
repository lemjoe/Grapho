package article

import (
	"fmt"
	"time"

	"github.com/lemjoe/md-blog/internal/models"
	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	q "github.com/ostafen/clover/v2/query"
)

type Article struct {
	collectionName string
	db             *c.DB
}
type articleSchema struct {
	FileName         string    `json:"file_name"`
	Title            string    `json:"article_title"`
	Author           string    `json:"article_author"`
	AuthorId         string    `json:"author_id"`
	CreationDate     time.Time `json:"creation_date"`
	ModificationDate time.Time `json:"modification_date"`
	IsLocked         bool      `json:"is_locked"`
	Id               string    `json:"id"`
}

func Init(db *c.DB) (*Article, error) {
	collection := Article{
		collectionName: "articles",
		db:             db,
	}
	err := db.CreateCollection(collection.collectionName)
	if err != nil {
		return nil, fmt.Errorf("unable to create collection[%s]: %w", collection.collectionName, err)
	}
	return &collection, nil
}
func (a *Article) CreateArticle(article models.Article) (models.Article, error) {
	doc := d.NewDocument()

	doc.Set("file_name", article.FileName)
	doc.Set("article_title", article.Title)
	doc.Set("article_author", article.Author)
	doc.Set("author_id", article.AuthorId)
	doc.Set("creation_date", time.Now())
	doc.Set("modification_date", time.Now())
	doc.Set("is_locked", article.IsLocked)

	docId, err := a.db.InsertOne(a.collectionName, doc)
	if err != nil {
		return models.Article{}, fmt.Errorf("unable to insert document[%s]: %w", a.collectionName, err)
	}
	return models.Article{
		FileName:         article.FileName,
		Title:            article.Title,
		Author:           article.Author,
		AuthorId:         article.AuthorId,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		IsLocked:         article.IsLocked,
		Id:               docId,
	}, nil
}
func (a *Article) GetAllArticles() ([]models.Article, error) {
	var articles []models.Article
	docs, err := a.db.FindAll(q.NewQuery(a.collectionName))
	if err != nil {
		return nil, fmt.Errorf("unable to find documents[%s]: %w", a.collectionName, err)
	}
	for _, doc := range docs {
		var article articleSchema
		err := doc.Unmarshal(&article)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal document[%s]: %w", a.collectionName, err)
		}

		articles = append(articles, models.Article{
			FileName:         article.FileName,
			Title:            article.Title,
			Author:           article.Author,
			AuthorId:         article.AuthorId,
			CreationDate:     article.CreationDate,
			ModificationDate: article.ModificationDate,
			IsLocked:         article.IsLocked,
			Id:               article.Id,
		})
	}
	return articles, nil
}

// GetArticleByFileName(fileName string) (models.Article, error)
func (a *Article) GetArticleByFileName(fileName string) (models.Article, error) {
	doc, err := a.db.FindFirst(q.NewQuery(a.collectionName).Where(q.Field("file_name").Eq(fileName)))
	if err != nil {
		return models.Article{}, fmt.Errorf("unable to find document[%s]: %w", a.collectionName, err)
	}
	var article articleSchema
	err = doc.Unmarshal(&article)
	if err != nil {
		return models.Article{}, fmt.Errorf("unable to unmarshal document[%s]: %w", a.collectionName, err)
	}
	return models.Article{
		FileName:         article.FileName,
		Title:            article.Title,
		Author:           article.Author,
		AuthorId:         article.AuthorId,
		CreationDate:     article.CreationDate,
		ModificationDate: article.ModificationDate,
		IsLocked:         article.IsLocked,
		Id:               article.Id,
	}, nil
}

// DeleteArticleByFileName(fileName string) error
func (a *Article) DeleteArticleByFileName(fileName string) error {
	err := a.db.Delete(q.NewQuery(a.collectionName).Where(q.Field("file_name").Eq(fileName)))
	if err != nil {
		return fmt.Errorf("unable to find document[%s]: %w", a.collectionName, err)
	}

	return nil
}

// UpdateArticleByFileName(fileName string) error
func (a *Article) UpdateArticleByFileName(fileName string) error {
	changes := make(map[string]interface{})
	changes["modification_date"] = time.Now()
	changes["is_locked"] = false
	err := a.db.Update(q.NewQuery(a.collectionName).Where(q.Field("file_name").Eq(fileName)), changes)
	return err
}
