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
	Title            string             `bson:"article_title"`
	Author           string             `bson:"article_author"`
	AuthorId         primitive.ObjectID `bson:"author_id"`
	CreationDate     time.Time          `bson:"creation_date"`
	ModificationDate time.Time          `bson:"modification_date"`
	IsLocked         bool               `bson:"is_locked"`
	Id               primitive.ObjectID `bson:"_id"`
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

func (a *Article) CreateArticle(article models.Article) (models.Article, error) {
	authorIdObj, err := primitive.ObjectIDFromHex(article.AuthorId)
	if err != nil {
		return models.Article{}, err
	}
	art := articleSchema{
		Title:            article.Title,
		Author:           article.Author,
		AuthorId:         authorIdObj,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
		IsLocked:         false,
	}

	res, err := a.ct.InsertOne(context.TODO(), bson.M{
		"article_title":     art.Title,
		"article_author":    art.Author,
		"author_id":         art.AuthorId,
		"creation_date":     art.CreationDate,
		"modification_date": art.ModificationDate,
		"is_locked":         art.IsLocked,
	})
	if err != nil {
		return models.Article{}, err
	}
	art.Id = res.InsertedID.(primitive.ObjectID)
	return models.Article{
		Title:            art.Title,
		Author:           art.Author,
		AuthorId:         art.AuthorId.Hex(),
		CreationDate:     art.CreationDate,
		ModificationDate: art.ModificationDate,
		IsLocked:         art.IsLocked,
		Id:               art.Id.Hex(),
	}, nil
}
func (a *Article) GetAllArticles() ([]models.Article, error) {
	var findedArticles []models.Article
	cur, err := a.ct.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var art articleSchema
		err := cur.Decode(&art)
		if err != nil {
			return nil, err
		}
		findedArticles = append(findedArticles, models.Article{
			Title:            art.Title,
			Author:           art.Author,
			AuthorId:         art.AuthorId.Hex(),
			CreationDate:     art.CreationDate,
			ModificationDate: art.ModificationDate,
			IsLocked:         art.IsLocked,
			Id:               art.Id.Hex(),
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedArticles, nil
}
func (a *Article) GetArticleById(id string) (models.Article, error) {
	artObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Article{}, err
	}
	var art articleSchema
	err = a.ct.FindOne(context.TODO(), bson.M{"_id": artObjId}).Decode(&art)
	if err != nil {
		return models.Article{}, err
	}
	return models.Article{
		Title:            art.Title,
		Author:           art.Author,
		AuthorId:         art.AuthorId.Hex(),
		CreationDate:     art.CreationDate,
		ModificationDate: art.ModificationDate,
		IsLocked:         art.IsLocked,
		Id:               art.Id.Hex(),
	}, nil
}
func (a *Article) DeleteArticleById(id string) error {
	artObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = a.ct.DeleteOne(context.TODO(), bson.M{"_id": artObjId})
	if err != nil {
		return err
	}

	return nil
}
func (a *Article) UpdateArticleById(id string) error {
	artObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = a.ct.UpdateOne(context.TODO(), bson.M{"_id": artObjId}, bson.M{"$set": bson.M{
		"modification_date": time.Now(),
		"is_locked":         false,
	}})
	if err != nil {
		return err
	}
	return nil
}
func (a *Article) LockArticleById(id string) error {
	artObjId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = a.ct.UpdateOne(context.TODO(), bson.M{"_id": artObjId}, bson.M{"$set": bson.M{
		"is_locked": true,
	}})
	if err != nil {
		return err
	}
	return nil
}
