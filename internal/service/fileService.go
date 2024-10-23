package service

import "os"

type fileService struct {
}

func NewFileService() *fileService {
	return &fileService{}
}

// ReadFile(fileName string) ([]byte, error)
func (f *fileService) ReadFile(fileName string) ([]byte, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// CreateNewFile(path string, body []byte) error
func (f *fileService) CreateNewFile(path string, body []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		return err
	}
	return nil
}

// WriteFile(path string, body []byte) error
func (f *fileService) WriteFile(path string, body []byte) error {
	err := os.WriteFile(path, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

// CreateFolder(path string) error
func (f *fileService) CreateFolder(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile(path string) error
func (f *fileService) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func (f *fileService) ListFolder(path string) ([][]byte, error) {
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var fileNames [][]byte
	for _, file := range list {
		fileNames = append(fileNames, []byte(file.Name()))
	}
	return fileNames, nil
}
