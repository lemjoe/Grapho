package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lemjoe/md-blog/internal/models"
	"github.com/lemjoe/md-blog/internal/repository/repotypes"
	"golang.org/x/crypto/bcrypt"
)

type migrationService struct {
	repository *repotypes.Repository
	artService *articleService
}

func NewMigrationService(repository *repotypes.Repository, artService *articleService) *migrationService {
	return &migrationService{
		repository: repository,
		artService: artService,
	}
}
func (m *migrationService) Migrate() error {
	if _, err := os.Stat("articles"); os.IsNotExist(err) {
		err = os.Mkdir("articles", os.ModePerm)
		if err != nil {
			return fmt.Errorf("migrate:\nunable to create articles folder: %w", err)
		}
	}
	_, err := m.repository.User.GetUserByUsername("admin")
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			password, _ := bcrypt.GenerateFromPassword([]byte("admin"), 10)
			newUsr, err := m.repository.User.CreateUser(models.User{
				UserName:     "admin",
				FullName:     "Administrator",
				Password:     string(password),
				Email:        "",
				IsAdmin:      true,
				Id:           "",
				LastLogin:    time.Now(),
				CreationDate: time.Now(),
			})
			fmt.Printf("migrate:\nadmin user created:[%+v]\n", newUsr)
			if err != nil {
				return fmt.Errorf("migrate:\nunable to create admin user: %w", err)
			}
		} else {
			return fmt.Errorf("migrate:\nunable to get admin user: %w", err)
		}
	}
	//fmt.Printf("migrate:\nadmin user:[%+v]\n", usr)
	mods, err := m.artService.GetArticlesList()
	if err != nil {
		if !strings.Contains(err.Error(), "unable to find documents") {
			return fmt.Errorf("migrate:\nunable to get articles list: %w", err)
		}
	}
	if len(mods) == 0 || (err != nil && strings.Contains(err.Error(), "unable to find documents")) {
		_, err := m.artService.CreateNewArticle("Welcome words", "admin", helloMessage)
		if err != nil {
			return fmt.Errorf("migrate:\nunable to create article: %w", err)
		}
	}

	return nil
}
