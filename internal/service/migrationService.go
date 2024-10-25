package service

import (
	"fmt"
	"os"
	"strings"
)

type migrationService struct {
	artService  *articleService
	userService *userService
}

func NewMigrationService(artService *articleService, userService *userService) *migrationService {
	return &migrationService{

		artService:  artService,
		userService: userService,
	}
}
func (m *migrationService) Migrate(adminPasswd string) error {
	if _, err := os.Stat("articles"); os.IsNotExist(err) {
		err = os.Mkdir("articles", os.ModePerm)
		if err != nil {
			return fmt.Errorf("migrate:\nunable to create articles folder: %w", err)
		}
	}

	_, err := m.userService.GetUserByName("admin")
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			newUsr, err := m.userService.CreateNewUser("admin", "Administrator", adminPasswd, "admin", true, true)
			fmt.Printf("migrate:\nadmin user created:[%+v]\n", newUsr)
			fmt.Println("please change the admin password first")
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
