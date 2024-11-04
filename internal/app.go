package internal

import (
	"github.com/lemjoe/Grapho/internal/config"
	"github.com/lemjoe/Grapho/internal/handler"
	"github.com/lemjoe/Grapho/internal/repository"
	"github.com/lemjoe/Grapho/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}
func (a *App) Run(version string) error {
	confDB, confApp, err := config.InitConfig("./.env")
	if err != nil {
		return err
	}
	service.InitLogs(confApp.MainLog)
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
	err = services.MigrationService.Migrate(confApp.AdminPasswd)
	if err != nil {
		return err
	}

	handlers := handler.NewHandler(services, bundle, version)
	err = handlers.Run(":" + confApp.Port)
	return err
}
