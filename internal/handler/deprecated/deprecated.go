package deprecated

import (
	"fmt"
	"net/http"

	"github.com/lemjoe/md-blog/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type handler struct {
	bundle   *i18n.Bundle
	services *service.Service
}

func Init(bundle *i18n.Bundle, services *service.Service) *handler {

	return &handler{
		bundle:   bundle,
		services: services,
	}
}

// Router
func (h *handler) Router() *http.ServeMux {
	fmt.Print("router\n")
	//http.HandleFunc("/ping", h.Ping)
	subrouter := http.NewServeMux()
	subrouter.HandleFunc("/show", h.ShowArticle)
	subrouter.HandleFunc("/edit", h.Editor)
	subrouter.HandleFunc("/delete", h.DeleteArticle)
	subrouter.HandleFunc("/add", h.UploadArticle)
	subrouter.HandleFunc("/upload", h.Upload)
	subrouter.HandleFunc("/download", h.DownloadArticle)
	subrouter.HandleFunc("/convert", h.MDConvert)
	subrouter.HandleFunc("/save", h.SaveFile)
	subrouter.HandleFunc("/singup", h.SingUp)
	subrouter.HandleFunc("/singin", h.SingIn)
	subrouter.HandleFunc("/", h.ArticleList)
	subrouter.HandleFunc("/ping", h.Ping)
	return subrouter
}

// Ping
func (h *handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("ping")
	w.Write([]byte("pong"))
}
