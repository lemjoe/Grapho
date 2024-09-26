package handler

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lemjoe/Grapho/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// SingUp
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))

	lang := curUser.Settings["language"]
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/sign-up.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	SingUpPageVars := models.PageVariables{ //store the date and time in a struct
		HomeButton: translation["homeButton"],
		Theme:      curUser.Settings["theme"],
	}
	err = t.Execute(w, SingUpPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                    // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func (h *Handler) SignUpPost(w http.ResponseWriter, r *http.Request) {
	log.Println("Registration form load")
	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")
	fullName := r.FormValue("fullname")

	_, err := h.services.UserService.GetUserById(login)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			newUsr, err := h.services.UserService.CreateNewUser(login, fullName, password, email, false)
			if err != nil {
				log.Print("unable to create user: ", err)
			} else {
				log.Printf("new user was created:[%+v]\n", newUsr)
			}
		} else {
			log.Print("unable to get user: ", err)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SingIn
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))

	lang := curUser.Settings["language"]
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/sign-in.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	SingInPageVars := models.PageVariables{ //store the date and time in a struct
		HomeButton: translation["homeButton"],
		Theme:      curUser.Settings["theme"],
	}
	err = t.Execute(w, SingInPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                    // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func (h *Handler) SignInPost(w http.ResponseWriter, r *http.Request) {
	log.Println("Login form load")
	login := r.FormValue("login")
	password := r.FormValue("password")

	user, err := h.services.UserService.GetUserByName(login)
	if err != nil {
		log.Print("invalid login or password: ", err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Print("invalid login or password: ", err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Print("failed to create token: ", err)
		return
	}
	log.Print("token created: ", tokenString)

	// Set cookie
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Path:     "",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	// Set cookie
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
