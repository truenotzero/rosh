package main

import (
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Note struct {
	Name    string
	Content string
}

func newNote(name, content string) Note {
	return Note{
		Name:    name,
		Content: content,
	}
}

type Notes = map[string]Note

type Page struct {
	Notes Notes
	Test  string
}

func newPage() Page {
	return Page{
		Notes: Notes{
			"lorem shmipsum": newNote("lorem shmipsum", "dolor sit amet"),
			"abcdef":         newNote("abcdef", "gefhijk"),
		},
		Test: "templating works",
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	e.Static("/static", "static")
	e.Static("/favicon.ico", "static/favicon.png")

	page := newPage()
	for _, n := range page.Notes {
		println(n.Name)
	}

	delay := func() {
		time.Sleep(1 * time.Second)
	}

	// alias for /notes
	e.GET("/", func(c echo.Context) error {
		delay()
		return c.Render(http.StatusOK, "index", page)
	})

	/// Get the names of all notes
	/// params - none
	/// values - none
	/// returns 200
	e.GET("/notes", func(c echo.Context) error {
		delay()
		return c.Render(http.StatusOK, "index", page)
	})

	/// Get the content of a note by name
	/// params - name: the note to get the content of
	/// values - none
	/// returns 200 & the note's content
	e.GET("/notes/:name", func(c echo.Context) error {
		delay()
		name := strings.TrimSpace(c.Param("name"))
		return c.String(http.StatusOK, page.Notes[name].Content)
	})

	/// Add a note, without overwriting an existing note
	/// params - none
	/// values - note-name: the name of the note to add/replace
	///			 note-data: the note's new data
	/// returns success: 200 & the new note element
	///			failure (note exists): 409 & error message
	e.PUT("/notes", func(c echo.Context) error {
		delay()
		name := strings.TrimSpace(c.FormValue("note-name"))
		data := c.FormValue("note-data")

		// check for presence of the key
		if _, b := page.Notes[name]; b {
			return c.String(http.StatusConflict, "note already exists")
		} else {
			note := newNote(name, data)
			page.Notes[name] = note
			return c.Render(http.StatusOK, "note", note)
		}
	})

	/// Update the content of an existing note
	/// params - none
	/// values - note-name: the name of the note to add/replace
	///			 note-data: the note's new data
	/// returns success: 200
	///			failure (note doesn't exists): 409 & error message
	e.PATCH("/notes", func(c echo.Context) error {
		delay()
		name := strings.TrimSpace(c.FormValue("note-name"))
		data := c.FormValue("note-data")

		// check for presence of the key
		if _, b := page.Notes[name]; b {
			note := newNote(name, data)
			page.Notes[name] = note
			return c.NoContent(http.StatusOK)
		} else {
			return c.String(http.StatusConflict, "note not found")
		}
	})

	/// Delete a note by name
	/// params - name: the name of the note to delete
	/// values - none
	/// returns 200
	e.DELETE("/notes/:name", func(c echo.Context) error {
		delay()
		name := strings.TrimSpace(c.Param("name"))

		delete(page.Notes, name)
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
