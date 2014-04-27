package cutecode

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
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

//takes appengine Context and returns the datastore key
func cutecodeKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "CuteCode", "default_cutecode", 0, nil)
}

func init() {
	http.HandleFunc("/", handler) //default handler
	http.HandleFunc("/signin", userHandler) //if route param is signin then invoke userHandler
	http.HandleFunc("/paste", codeHandler) //if route param is paste then invoke codeHandler
	http.HandleFunc("/save", saveHandler) //if route param is save then invoke saveHandler
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("CuteCode").Ancestor(cutecodeKey(c)).Order("-CreatedAt").Limit(10)

	codes := make([]Code, 0, 10)
	if _, err := q.GetAll(c, &codes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, codes)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("templates/signin.html")
	t.Execute(w, nil)
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/codeform.html")
	t.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	t, _ := template.ParseFiles("templates/showcode.html")
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
	t.Execute(w, &code)
}
