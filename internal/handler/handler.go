package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lemjoe/md-blog/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Handler struct {
	services *service.Service
	bundle   *i18n.Bundle
}

func NewHandler(services *service.Service, bundle *i18n.Bundle) *Handler {
	return &Handler{
		services: services,
		bundle:   bundle,
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
	r.Handle("/singup", http.HandlerFunc(h.SingUp))
	r.Handle("/singin", http.HandlerFunc(h.SingIn))

	// import resources
	r.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(http.Dir("./lib/"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))

	log.Print("Server is running on port ", port)
	return http.ListenAndServe(port, r)
}
