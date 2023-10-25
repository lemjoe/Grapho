package repository

import (
	"fmt"

	"github.com/lemjoe/md-blog/internal/models"
	"github.com/lemjoe/md-blog/internal/repository/cloverdb"
	"github.com/lemjoe/md-blog/internal/repository/repotypes"
)

type DB interface {
	Close()
	NewRepository() (*repotypes.Repository, error)
}

func InitializeDB(dbType string, conf models.ConfigDB) (DB, error) {
	switch dbType {
	case "cloverdb":
		return cloverdb.ConnectDB(conf.Path)
	default:
		return nil, fmt.Errorf("invalid db type: %s", dbType)
	}
}
