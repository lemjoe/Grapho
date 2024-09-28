// Users management
package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/lemjoe/Grapho/internal/models"
)

func (h *Handler) GetUsersList(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))

	log.Println("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		log.Println("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin {
		log.Println("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	theme := curUser.Settings["theme"]
	lang := curUser.Settings["language"]
	translation := Localizer([]string{"listOfArticles", "homeButton", "addButton", "lastModification"}, lang, h.bundle)
	// log.Println(translation)

	users, err := h.services.UserService.GetUsersList()
	if err != nil {
		log.Println(err)
	}

	html := "<h1>" + "List of registered users:" + "</h1><table><tr><th>#</th><th>Name</th><th>Full Name</th><th>Email</th><th>Admin?</th><th>Manage</th></tr>"
	editImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/" + theme + "/edit-pen.png\" alt=\"Edit\" title=\"Edit\">"
	// deleteImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/" + theme + "/red-trash-can.png\" alt=\"Edit\" title=\"Edit\">"

	// if len(docs) == 0 {
	// 	html += "<p>There is no articles here! Why don't you add one?"
	// }
	for i, user := range users {
		log.Println(user)
		html += "<tr><td>" + strconv.Itoa(i+1) + "</td><td>" + user.UserName + "</td><td>" + user.FullName + "</td><td>" + user.Email + "</td><td>" + strconv.FormatBool(user.IsAdmin) + "</td><td></a><a href='manage_user?md=" + user.Id + "'><i>" + editImg + "</i></a></td></tr>"
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
		log.Print("Template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("Template executing error: ", err) //log it
	}

}
