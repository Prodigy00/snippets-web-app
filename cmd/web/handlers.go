package main

import (
	"errors"
	"fmt"
	"github.com/Prodigy00/snippetbox/internal/models"
	"net/http"
	"strconv"
)

// eq to func logicFunc(res, req){} in express
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(queryParam)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	//w.Write([]byte("Display a specific snippet"))
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST") //needs to happen before you call write/writeHead!
		//w.WriteHeader(405)
		//w.Write([]byte("Method Not Allowed"))
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Que sera"
	content := "Whatever will be,\nwill be!\n\n- Anonymous!"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}

	redirectUrl := fmt.Sprintf("/snippet/view?id=%d", id)

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
