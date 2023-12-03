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
	userService := NewUserService(repository)
	migrateService := NewMigrationService(artService, userService)

	return &Service{
		repository:       repository,
		FileService:      fileServiceInstance,
		ArticleService:   artService,
		MigrationService: migrateService,
		UserService:      userService,
	}

}
