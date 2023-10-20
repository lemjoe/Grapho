package user

import (
	"fmt"

	c "github.com/ostafen/clover/v2"
)

type User struct {
	collectionName string
	db             *c.DB
}

func Init(db *c.DB) (*User, error) {
	collection := User{
		collectionName: "users",
		db:             db,
	}
	err := db.CreateCollection(collection.collectionName)
	if err != nil {
		return nil, fmt.Errorf("unable to create collection[%s]: %w", collection.collectionName, err)
	}
	return &collection, nil
}
