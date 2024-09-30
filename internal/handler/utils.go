package handler

import (
	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func MdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.Footnotes | parser.MathJax | parser.DefinitionLists | parser.Titleblock | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.FootnoteReturnLinks
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func (h *Handler) GetCurrentUser(userID string) *models.User {
	logger := service.GetLogger()
	curUser, err := h.services.UserService.GetUserById(userID)
	if err != nil {
		logger.Error(err)
		return &models.User{
			UserName: "guest",
			FullName: "Guest",
			IsAdmin:  false,
			Settings: service.DefaultUserSettings,
		}
	}
	return &curUser
}

func Localizer(input []string, lang string, bundle *i18n.Bundle) map[string]string {
	// defaultLang := "en"
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("lang/active.ru.toml")
	localizer := i18n.NewLocalizer(bundle, lang)
	localization := make(map[string]string)
	output := make(map[string]string)

	// Translation strings
	localization["listOfArticles"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ListOfArticles",
			Other: "List of available articles:",
		},
	})
	localization["homeButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HomeButton",
			Other: "Back to home page",
		},
	})
	localization["addButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "AddButton",
			Other: "Add an article",
		},
	})
	localization["lastModification"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "LastModification",
			Other: "Last modification",
		},
	})
	localization["pageTitle"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Title",
			Other: "Articles list",
		},
	})
	localization["user"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "User",
			Other: "User",
		},
	})
	localization["register"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Register",
			Other: "register",
		},
	})
	localization["login"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Login",
			Other: "login",
		},
	})
	localization["logout"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Logout",
			Other: "logout",
		},
	})

	for _, val := range input {
		output[val] = localization[val]
	}

	return output
}
