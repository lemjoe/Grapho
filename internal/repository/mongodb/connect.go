package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/lemjoe/Grapho/internal/repository/mongodb/article"
	"github.com/lemjoe/Grapho/internal/repository/mongodb/user"
	"github.com/lemjoe/Grapho/internal/repository/repotypes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
}
type DB struct {
	Driver      *mongo.Database
	Collections *Collection
	client      *mongo.Client
}

func ConnectDB(url, dbname string, user string, password string) (*DB, error) {
	if user == "" {
		url = "mongodb://" + url
	} else {
		url = "mongodb://" + user + ":" + password + "@" + url
	}
	fmt.Printf("url: %s, dbname: %s\n", url, dbname)
	//var collection *mongo.Collection
	var ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	databaseInstance := &DB{
		Driver: client.Database(dbname),
	}

	return databaseInstance, nil

}

func (db *DB) Close() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (db *DB) NewRepository() (*repotypes.Repository, error) {
	user, err := user.Init(db.Driver)
	if err != nil {
		return nil, err
	}
	art, err := article.Init(db.Driver)
	if err != nil {
		return nil, err
	}
	return &repotypes.Repository{
		User:    user,
		Article: art,
	}, nil
}
