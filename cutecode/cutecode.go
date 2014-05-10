package cutecode

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"net/http"
	"time"
)

type Code struct {
	Title     string
	Content   string
	UrlKey    string
	Lang      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var indexTemplate = template.Must(template.ParseFiles("templates/layout.html", "templates/index.html"))
var showcodeTemplate = template.Must(template.ParseFiles("templates/layout.html", "templates/showcode.html"))

//takes appengine Context and returns the datastore key
func cutecodeKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "CuteCode", "default_cutecode", 0, nil)
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/show", showHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := indexTemplate.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	code := Code{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Lang:    r.FormValue("lang"),
	}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "CuteCode", cutecodeKey(c)), &code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templateerr := showcodeTemplate.ExecuteTemplate(w, "layout", &code)
	if templateerr != nil {
		http.Error(w, templateerr.Error(), http.StatusInternalServerError)
	}
}

func showHandler(w http.ResponseWriter, r *http.Request) {
    code := Code{
		Title:   "Hello",
		Content: "Bello",
		Lang:    "Go",
    }

    templateerr := showcodeTemplate.ExecuteTemplate(w, "layout", &code)
    if templateerr != nil {
        http.Error(w, templateerr.Error(), http.StatusInternalServerError)
    }
}
