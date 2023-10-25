package service

import "github.com/lemjoe/md-blog/internal/repository"

type Service struct {
	repository       *repository.Repository
	FileService      FileService
	ArticleService   ArticleService
	MigrationService MigrationService
}

func NewService(repository *repository.Repository) *Service {
	fileServiceInstance := NewFileService()
	artService := NewArticleService(repository, fileServiceInstance)
	migrateService := NewMigrationService(repository, artService)
	return &Service{
		repository:       repository,
		FileService:      fileServiceInstance,
		ArticleService:   artService,
		MigrationService: migrateService,
	}

}
