package repository

import (
	"github.com/lemjoe/md-blog/internal/repository/cloverdb/article"
	"github.com/lemjoe/md-blog/internal/repository/cloverdb/user"
	"github.com/lemjoe/md-blog/internal/repository/repotypes"
	"github.com/ostafen/clover/v2"
)

func NewRepository(db *clover.DB) (*repotypes.Repository, error) {
	user, err := user.Init(db)
	if err != nil {
		return nil, err
	}
	article, err := article.Init(db)
	if err != nil {
		return nil, err
	}
	return &repotypes.Repository{
		User:    user,
		Article: article,
	}, nil
}
