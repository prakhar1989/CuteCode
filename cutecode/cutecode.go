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
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	c := appengine.NewContext(r)
	q := datastore.NewQuery("CuteCode").Ancestor(cutecodeKey(c)).Order("-CreatedAt").Limit(10)

	codes := make([]Code, 0, 10)
	if _, err := q.GetAll(c, &codes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := indexTemplate.ExecuteTemplate(w, "layout", codes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	code := Code{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}

	key := datastore.NewIncompleteKey(c, "CuteCode", cutecodeKey(c))
	_, err := datastore.Put(c, key, &code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templerr := showcodeTemplate.ExecuteTemplate(w, "layout", &code)
	if templerr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
