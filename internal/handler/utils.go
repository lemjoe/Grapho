package handler

import (
	"log"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/lemjoe/md-blog/internal/models"
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
	curUser, err := h.services.UserService.GetUserById(userID)
	if err != nil {
		log.Println(err)
		return &models.User{
			UserName: "guest",
			FullName: "Guest",
			IsAdmin:  false,
		}
	}
	return &curUser
}

func Localizer(input []string, lang string, bundle *i18n.Bundle) map[string]string {
	defaultLang := "en"
	localizer := i18n.NewLocalizer(bundle, lang, defaultLang)
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

	for _, val := range input {
		output[val] = localization[val]
	}

	return output
}
