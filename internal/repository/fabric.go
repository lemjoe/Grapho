package repository

import (
	"fmt"

	"github.com/lemjoe/md-blog/internal/repository/cloverdb"
	"github.com/lemjoe/md-blog/internal/repository/repotypes"
)

type DB interface {
	Close()
	NewRepository() (*repotypes.Repository, error)
}
type ConfigDB struct {
	Path string
}

func InitializeDB(dbType string, conf ConfigDB) (DB, error) {
	switch dbType {
	case "cloverdb":
		return cloverdb.ConnectDB(conf.Path)
	default:
		return nil, fmt.Errorf("invalid db type: %s", dbType)
	}
}
