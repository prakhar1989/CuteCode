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

//caching of templates
var indexTemplate = template.Must(template.ParseFiles("templates/layout.html","templates/index.html"))
var signinTemplate = template.Must(template.ParseFiles("templates/layout.html","templates/signin.html"))
var codeformTemplate = template.Must(template.ParseFiles("templates/layout.html","templates/codeform.html"))
var showcodeTemplate = template.Must(template.ParseFiles("templates/layout.html","templates/showcode.html"))

//takes appengine Context and returns the datastore key
func cutecodeKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "CuteCode", "default_cutecode", 0, nil)
}

func init() {
	http.HandleFunc("/", handler)           //default handler
	http.HandleFunc("/signin", userHandler) //if route param is signin then invoke userHandler
	http.HandleFunc("/paste", codeHandler)  //if route param is paste then invoke codeHandler
	http.HandleFunc("/save", saveHandler)   //if route param is save then invoke saveHandler
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
	err:=indexTemplate.ExecuteTemplate(w,"layout", codes) //this mention of layout stems from the {{layout}} of the layout page. for more http://blog.xcai.net/golang/templates
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
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
	err := signinTemplate.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	err := codeformTemplate.ExecuteTemplate(w, "layout", nil)
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
	templerr := showcodeTemplate.ExecuteTemplate(w,"layout", &code)
	if templerr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
