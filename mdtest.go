package main

import (
	"github.com/gomarkdown/markdown"

	"html/template"
	"log"
	"net/http"
	"os"
)

type PageVariables struct {
	Md        string
	MDArticle template.HTML
}

var toTheTop = []byte("\n<a href=\"#top\"><i>back to top</i></a>")

func init() {
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	// http.Handle("/articles/", http.StripPrefix("/articles/", http.FileServer(http.Dir("articles"))))
}

func main() {
	http.HandleFunc("/show", ShowArticle)
	http.HandleFunc("/edit", Editor)
	http.HandleFunc("/", ArticleList)
	http.HandleFunc("/callme", CallMe)
	log.Print("Server is running on port 3608")
	log.Fatal(http.ListenAndServe(":3608", nil))
}

func ShowArticle(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	md, err := os.ReadFile("articles/" + artclPath) // just pass the file name
	if err != nil {
		log.Print("MD file open error: ", err, artclPath)
	}
	// always normalize newlines!
	md = markdown.NormalizeNewlines(md)
	html := append(markdown.ToHTML(md, nil, nil), toTheTop[:]...)

	HomePageVars := PageVariables{ //store the date and time in a struct
		MDArticle: template.HTML(html),
	}

	t, err := template.ParseFiles("template.html") //parse the html file homepage.html
	if err != nil {                                // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func Editor(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	md, err := os.ReadFile("articles/" + artclPath) // just pass the file name
	if err != nil {
		log.Print("MD file open error: ", err)
	}
	// always normalize newlines!
	md = markdown.NormalizeNewlines(md)
	html := markdown.ToHTML(md, nil, nil)

	HomePageVars := PageVariables{ //store the date and time in a struct
		Md:        string(md),
		MDArticle: template.HTML(html),
	}

	t, err := template.ParseFiles("markdown_editor.html") //parse the html file homepage.html
	if err != nil {                                       // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func ArticleList(w http.ResponseWriter, r *http.Request) {

	f, err := os.Open("articles")
	if err != nil {
		log.Print("Articles directory open error: ", err)
		return
	}

	files, err := f.Readdir(0)
	if err != nil {
		log.Print("Articles read error: ", err)
		return
	}

	html := "<h2>List of available articles:</h2><ul>"

	for _, v := range files {
		if !v.IsDir() {
			html += "<li>" + "<a href='show?md=" + v.Name() + "'>" + v.Name() + "</a><i> (Last modification: " + v.ModTime().Format("2006-Jan-02") + ") </i><a href='edit?md=" + v.Name() + "'><i>edit</i></a></li>"
		}
	}

	html += "</ul>"

	HomePageVars := PageVariables{ //store the date and time in a struct
		MDArticle: template.HTML(html),
	}

	t, err := template.ParseFiles("template.html") //parse the html file homepage.html
	if err != nil {                                // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func CallMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You called me!"))
}
