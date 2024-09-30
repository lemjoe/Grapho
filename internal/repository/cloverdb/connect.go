package cloverdb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lemjoe/Grapho/internal/repository/cloverdb/article"
	"github.com/lemjoe/Grapho/internal/repository/cloverdb/user"
	"github.com/lemjoe/Grapho/internal/repository/repotypes"
	"github.com/ostafen/clover/v2"
)

type DB struct {
	DB *clover.DB
}

func ConnectDB(dir string) (*DB, error) {
	//relative path to absolute
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	fmt.Printf("db path: %s\n", absDir)
	//check if dir exists
	_, err = os.Stat(absDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(absDir, 0740)
			if err != nil {
				return nil, fmt.Errorf("unable to create database: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unable to create database: %w", err)
		}
	}

	db, err := clover.Open(absDir)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB: db,
	}, nil
}
func (d *DB) Close() {
	d.DB.Close()
}

// NewRepository() (*Repository, error)
func (d *DB) NewRepository() (*repotypes.Repository, error) {
	user, err := user.Init(d.DB)
	if err != nil {
		return nil, err
	}
	article, err := article.Init(d.DB)
	if err != nil {
		return nil, err
	}
	return &repotypes.Repository{
		User:    user,
		Article: article,
	}, nil
}
