package handler

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func MdToHTML(md []byte, editMode bool) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.Footnotes | parser.MathJax | parser.DefinitionLists | parser.Titleblock | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.FootnoteReturnLinks
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	result := markdown.Render(doc, renderer)
	if !editMode {
		copyCode := []byte("<pre><a href=\"\" class=\"copy-code\" title=\"Copy to clipboard\"><img style=\"padding:0px;opacity:0.6;filter:alpha(opacity=60);\" width=\"16\" height=\"16\" src=\"../images/copy.png\"></a><code")
		result = bytes.ReplaceAll(result, []byte("<pre><code"), copyCode)
	}
	return result
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

	// Top menu
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
	localization["editButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "EditButton",
			Other: "Edit",
		},
	})
	localization["deleteButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "DeleteButton",
			Other: "Delete",
		},
	})
	localization["switchModeButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SwitchModeButton",
			Other: "Switch mode",
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
	localization["download"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Download",
			Other: "Download",
		},
	})

	// Page titles
	localization["titleMain"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleMain",
			Other: "Articles list",
		},
	})
	localization["titleEdit"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleEdit",
			Other: "Edit article",
		},
	})
	localization["titleUploadArt"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleUploadArt",
			Other: "Upload an article",
		},
	})
	localization["titleSignUp"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleSignUp",
			Other: "Sign up form",
		},
	})
	localization["titleLogin"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleLogin",
			Other: "Log in form",
		},
	})
	localization["titleUserSettings"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleUserSettings",
			Other: "User settings",
		},
	})
	localization["titleAdmUsersList"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleAdmUsersList",
			Other: "Admin panel | Users list",
		},
	})
	localization["titleAdmManageUser"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleAdmManageUser",
			Other: "Admin panel | Manage user",
		},
	})
	localization["titleLicenses"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleLicenses",
			Other: "License information",
		},
	})

	// Strings
	localization["listOfArticles"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ListOfArticles",
			Other: "List of available articles:",
		},
	})
	localization["lastModification"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "LastModification",
			Other: "Last modification",
		},
	})
	localization["uploadedBy"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "UploadedBy",
			Other: "Uploaded by",
		},
	})
	localization["backToTop"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BackToTop",
			Other: "back to top",
		},
	})
	localization["preview"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Preview",
			Other: "Preview",
		},
	})
	localization["save"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "save",
			Other: "Save",
		},
	})
	localization["articleTitle"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ArticleTitle",
			Other: "Article Title",
		},
	})
	localization["upload"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Upload",
			Other: "Upload",
		},
	})
	localization["by"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "By",
			Other: "by",
		},
	})
	localization["loginString"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "LoginString",
			Other: "login",
		},
	})
	localization["passwordString"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PasswordString",
			Other: "password",
		},
	})
	localization["fullNameString"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "FullNameString",
			Other: "full name",
		},
	})
	localization["onlyLatin"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "OnlyLatin",
			Other: "Only latin downcase letters and numbers",
		},
	})
	localization["mustContain"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "MustContain",
			Other: "Must contain upcase, downcase, number and at least 8 characters",
		},
	})
	localization["mustBeEmail"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "MustBeEmail",
			Other: "Must be a valid e-mail address",
		},
	})
	localization["submitButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SubmitButton",
			Other: "Submit",
		},
	})
	localization["loginButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "loginButton",
			Other: "Login",
		},
	})
	localization["settings"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Settings",
			Other: "Settings",
		},
	})
	localization["security"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Security",
			Other: "Security",
		},
	})
	localization["currentPass"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "CurrentPass",
			Other: "current password",
		},
	})
	localization["newPass"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "NewPass",
			Other: "new password",
		},
	})
	localization["reNewPass"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ReNewPass",
			Other: "retype new password",
		},
	})
	localization["changePassButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ChangePassButton",
			Other: "Change password",
		},
	})
	localization["misc"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Misc",
			Other: "Misc",
		},
	})
	localization["language"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Language",
			Other: "language",
		},
	})
	localization["colorTheme"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ColorTheme",
			Other: "color theme",
		},
	})
	localization["changeSettingsButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ChangeSettingsButton",
			Other: "Change settings",
		},
	})
	localization["adminPanel"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "AdminPanel",
			Other: "Admin Panel",
		},
	})
	localization["managePortalUsers"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ManagePortalUsers",
			Other: "manage portal users",
		},
	})
	localization["listOfUsers"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ListOfUsers",
			Other: "List of registered users",
		},
	})
	localization["thName"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ThName",
			Other: "Name",
		},
	})
	localization["thFullName"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ThFullName",
			Other: "Full name",
		},
	})
	localization["thManage"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ThManage",
			Other: "Manage",
		},
	})
	localization["manageUser"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ManageUser",
			Other: "Manage user",
		},
	})
	localization["userStr"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "UserStr",
			Other: "user",
		},
	})
	localization["isAdmin"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "IsAdmin",
			Other: "is admin?",
		},
	})
	localization["thWriter"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Writer",
			Other: "Writer?",
		},
	})
	localization["isWriter"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "IsWriter",
			Other: "is writer?",
		},
	})
	localization["changeUserButton"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ChangeUserButton",
			Other: "Change user",
		},
	})
	localization["backToUsersList"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BackToUsersList",
			Other: "back to users list",
		},
	})
	localization["noArticles"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "NoArticles",
			Other: "There is no articles here! Why don't you add one?",
		},
	})
	localization["licenseInformation"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "LicenseInformation",
			Other: "License information",
		},
	})
	localization["binariesAndSC"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BinariesAndSC",
			Other: "Binaries and source code of this software are available",
		},
	})
	localization["here"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Here",
			Other: "here",
		},
	})
	localization["underTerms"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "UnderTerms",
			Other: "under the terms of BSD 3-Clause License.",
		},
	})
	localization["thirdPartyLicenses"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ThirdPartyLicenses",
			Other: "Information about the licenses of the software used in the development of Grapho. The names of developers and organizations are not mentioned for the purpose of endorsing or promoting Grapho.",
		},
	})
	localization["uploadAnArticle"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "UploadAnArticle",
			Other: "Upload an article",
		},
	})
	localization["fromScratch"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "FromScratch",
			Other: "or write it from scratch",
		},
	})

	for _, val := range input {
		output[val] = localization[val]
	}

	return output
}

// Show licinses information
func (h *Handler) ShowLicenses(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	lang := curUser.Settings["language"]
	translation := Localizer(localization, lang, h.bundle)

	md, err := h.services.FileService.ReadFile("LICENSE")
	if err != nil {
		logger.Error("License file open error: ", err)
	}

	licensesPath, err := h.services.FileService.ListFolder("lib/licenses")
	if err != nil {
		logger.Error("Can't get licenses list: ", err)
	}
	var licenses []string
	for _, license := range licensesPath {
		str := strings.Replace(string(license), "^", "/", 1)
		licenses = append(licenses, str)
	}

	// always normalize newlines!
	html := MdToHTML(md, false)

	ArticlePageVars := models.PageVariables{ //store the date and time in a struct
		MDArticle:    template.HTML(html),
		UserName:     curUser.FullName,
		Theme:        curUser.Settings["theme"],
		Translation:  translation,
		Licenses:     licenses,
		LicensesPath: licensesPath,
		Title:        translation["titleLicenses"],
	}

	t, err := template.ParseFiles("lib/templates/licenses.html") //parse the html file homepage.html
	if err != nil {                                              // if there is an error
		logger.Error("template parsing error: ", err) // log it
	}
	err = t.Execute(w, ArticlePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		logger.Error("template executing error: ", err) //log it
	}
}
