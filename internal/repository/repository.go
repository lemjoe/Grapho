package repository

import (
	"github.com/lemjoe/md-blog/internal/repository/cloverdb/user"
	"github.com/ostafen/clover/v2"
)

type User interface {
}
type Repository struct {
	User User
}

func NewRepository(db *clover.DB) (*Repository, error) {
	user, err := user.Init(db)
	if err != nil {
		return nil, err
	}
	return &Repository{
		User: user,
	}, nil
}
