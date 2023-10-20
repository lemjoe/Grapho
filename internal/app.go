package internal

import (
	"github.com/lemjoe/md-blog/internal/handler"
	"github.com/lemjoe/md-blog/internal/repository"
	"github.com/lemjoe/md-blog/internal/repository/cloverdb"
	"github.com/lemjoe/md-blog/internal/service"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}
func (a *App) Run() error {
	db, err := cloverdb.ConnectDB("db")
	if err != nil {
		return err
	}
	defer db.Close()
	repos, err := repository.NewRepository(db.DB)
	if err != nil {
		return err
	}
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	err = handlers.Run(":4007")
	return err
}
