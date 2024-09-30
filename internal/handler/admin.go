// Users management
package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
)

func (h *Handler) GetUsersList(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin {
		logger.Error("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	theme := curUser.Settings["theme"]
	lang := curUser.Settings["language"]
	translation := Localizer([]string{"listOfArticles", "homeButton", "addButton", "lastModification"}, lang, h.bundle)

	users, err := h.services.UserService.GetUsersList()
	if err != nil {
		logger.Error(err)
	}

	html := "<h1>" + "List of registered users:" + "</h1><table><tr><th>#</th><th>Name</th><th>Full Name</th><th>Email</th><th>Admin?</th><th>Manage</th></tr>"
	editImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/" + theme + "/edit-pen.png\" alt=\"Edit\" title=\"Edit\">"
	// deleteImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/" + theme + "/red-trash-can.png\" alt=\"Edit\" title=\"Edit\">"

	// if len(docs) == 0 {
	// 	html += "<p>There is no articles here! Why don't you add one?"
	// }
	for i, user := range users {
		logger.Info(user)
		html += "<tr><td>" + strconv.Itoa(i+1) + "</td><td>" + user.UserName + "</td><td>" + user.FullName + "</td><td>" + user.Email + "</td><td>" + strconv.FormatBool(user.IsAdmin) + "</td><td></a><a href='manageuser?usr=" + user.Id + "'><i>" + editImg + "</i></a></td></tr>"
	}

	html += "</table>"

	HomePageVars := models.PageVariables{ //store the date and time in a struct
		MDArticle:  template.HTML(html),
		HomeButton: translation["homeButton"],
		AddButton:  translation["addButton"],
		Title:      "Admin | Users list",
		UserName:   curUser.FullName,
		Theme:      theme,
	}

	t, err := template.ParseFiles("lib/templates/home.html") //parse the html file homepage.html
	if err != nil {                                          // if there is an error
		logger.Error("Template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		logger.Error("Template executing error: ", err) //log it
	}

}

// Show article
func (h *Handler) ManageUser(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin {
		logger.Error("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	lang := curUser.Settings["language"]
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	usrId := r.URL.Query().Get("usr")
	managedUsr, err := h.services.UserService.GetUserById(usrId)
	if err != nil {
		logger.Error("Can't load user data: ", err, managedUsr)
	}

	mngUsrIsAdm := ""
	if managedUsr.IsAdmin {
		mngUsrIsAdm = "checked"
	}

	ManageUserPageVars := models.PageVariables{ //store the date and time in a struct
		Title:               "Admin | Manage user",
		HomeButton:          translation["homeButton"],
		UserName:            curUser.FullName,
		Theme:               curUser.Settings["theme"],
		ManagedUserName:     managedUsr.UserName,
		ManagedUserFullName: managedUsr.FullName,
		ManagedUserEmail:    managedUsr.Email,
		ManagedUserIsAdmin:  mngUsrIsAdm,
	}

	t, err := template.ParseFiles("lib/templates/manage-user.html") //parse the html file homepage.html
	if err != nil {                                                 // if there is an error
		logger.Error("template parsing error: ", err) // log it
	}
	err = t.Execute(w, ManageUserPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                        // if there is an error
		logger.Error("template executing error: ", err) //log it
	}
}

func (h *Handler) ChangeUser(w http.ResponseWriter, r *http.Request) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin {
		logger.Error("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	user_name := r.FormValue("user_name")
	full_name := r.FormValue("full_name")
	e_mail := r.FormValue("e_mail")
	is_admin := false
	if r.FormValue("is_admin") == "on" {
		is_admin = true
	}

	managedUser, err := h.services.UserService.GetUserByName(user_name)
	if err != nil {
		logger.Error("Can't load user data: ", err)
	}

	err = h.services.UserService.UpdateUserData(managedUser.Id, full_name, e_mail, is_admin)
	if err != nil {
		logger.Error("Something went wrong: ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
