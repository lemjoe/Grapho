package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lemjoe/Grapho/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Handler struct {
	services *service.Service
	bundle   *i18n.Bundle
	version  string
}

func NewHandler(services *service.Service, bundle *i18n.Bundle, version string) *Handler {
	return &Handler{
		services: services,
		bundle:   bundle,
		version:  version,
	}
}

func (h *Handler) Run(port string) error {
	r := mux.NewRouter()
	r.Handle("/", http.HandlerFunc(h.GetArticlesList)).Methods("GET")
	r.Handle("/show", http.HandlerFunc(h.ShowArticle)).Methods("GET")
	r.Handle("/edit", http.HandlerFunc(h.Editor)).Methods("GET")
	r.Handle("/delete", http.HandlerFunc(h.DeleteArticle)).Methods("GET")
	r.Handle("/add", http.HandlerFunc(h.UploadArticle))
	r.Handle("/upload", http.HandlerFunc(h.Upload))
	r.Handle("/download", http.HandlerFunc(h.DownloadArticle))
	r.Handle("/convert", http.HandlerFunc(h.MDConvert))
	r.Handle("/save", http.HandlerFunc(h.SaveFile))
	r.Handle("/signup", http.HandlerFunc(h.SignUp)).Methods("GET")
	r.Handle("/signin", http.HandlerFunc(h.SignIn)).Methods("GET")
	r.Handle("/signup", http.HandlerFunc(h.SignUpPost)).Methods("POST")
	r.Handle("/signin", http.HandlerFunc(h.SignInPost)).Methods("POST")
	r.Handle("/logout", http.HandlerFunc(h.Logout))
	r.Handle("/settings", http.HandlerFunc(h.Settings)).Methods("GET")
	r.Handle("/changepasswd", http.HandlerFunc(h.ChangeUserPassword)).Methods("POST")
	r.Handle("/changetheme", http.HandlerFunc(h.ChangeTheme))
	r.Handle("/changeusersettings", http.HandlerFunc(h.ChangeSettings))
	r.Handle("/admin", http.HandlerFunc(h.GetUsersList)).Methods("GET")
	r.Handle("/manageuser", http.HandlerFunc(h.ManageUser)).Methods("GET")
	r.Handle("/changeuser", http.HandlerFunc(h.ChangeUser)).Methods("POST")
	r.Handle("/licenses", http.HandlerFunc(h.ShowLicenses)).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(h.PageNotFound)
	r.Use(authMiddleware) // JWT check

	// import resources
	r.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(http.Dir("./lib/"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))

	logger := service.GetLogger()
	logger.Info("Server is running on port ", port)
	return http.ListenAndServe(port, r)
}
