package service

import "github.com/lemjoe/md-blog/internal/repository/repotypes"

type Service struct {
	repository       *repotypes.Repository
	FileService      FileService
	ArticleService   ArticleService
	MigrationService MigrationService
	UserService      UserService
}

func NewService(repository *repotypes.Repository) *Service {
	fileServiceInstance := NewFileService()
	artService := NewArticleService(repository, fileServiceInstance)
	migrateService := NewMigrationService(repository, artService)
	userService := NewUserService(repository)
	return &Service{
		repository:       repository,
		FileService:      fileServiceInstance,
		ArticleService:   artService,
		MigrationService: migrateService,
		UserService:      userService,
	}

}
