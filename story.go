package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

// init function
// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized.
// Besides initializations that cannot be expressed as declarations, a common use
// of init functions is to verify or repair correctness of the program state
// before real execution begins.

// template.Must()
// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil.
// create template with New("name"), name could be empty ""
// then Parse the given text for the template, here it's html
// later will Execute the template on chapters
func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
        <p>{{.}}</p>
    {{end}}
    <ul>
    {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
    </ul>
</body>
</html>`

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// JSONStory takes the given json file and turns it into a Story map
func JSONStory(r io.Reader) (Story, error) {
	// json.NewDecoder() --> pass in an io.Reader
	// json.Marshal() --> pass in a byte slice
	dd := json.NewDecoder(r)
	var story Story

	if err := dd.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// NewHandler takes in the Story map created by JSONStory
// and returns a new http.Handler interface
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// trim space isn't super necessary but a good check
	path := strings.TrimSpace(r.URL.Path)
	// if there's no path, start from the beginning of the story
	if path == "" || path == "/" {
		path = "/intro"
	}
	// get rid of leading slash
	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		// use the template created in init to execute
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
