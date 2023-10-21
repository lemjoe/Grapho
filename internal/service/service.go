package service

import "github.com/lemjoe/md-blog/internal/repository"

type Service struct {
	repository *repository.Repository
	FileService
}

func NewService(repository *repository.Repository) *Service {
	fileServiceInstance := NewFileService()
	return &Service{
		repository:  repository,
		FileService: fileServiceInstance,
	}
}
