package internal

import (
	"github.com/lemjoe/md-blog/internal/config"
	"github.com/lemjoe/md-blog/internal/handler"
	"github.com/lemjoe/md-blog/internal/repository"
	"github.com/lemjoe/md-blog/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}
func (a *App) Run() error {
	confDB, confApp, err := config.InitConfig("./.env")
	if err != nil {
		return err
	}
	db, err := repository.InitializeDB(confDB.DbType, confDB)
	if err != nil {
		return err
	}
	repos, err := db.NewRepository()
	if err != nil {
		return err
	}
	bundle := i18n.NewBundle(language.English)
	services := service.NewService(repos)
	err = services.MigrationService.Migrate()
	if err != nil {
		return err
	}

	handlers := handler.NewHandler(services, bundle)
	err = handlers.Run(":" + confApp.Port)
	return err
}
