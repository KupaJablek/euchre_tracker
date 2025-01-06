package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/KupaJablek/euchre_tracker/internal"
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

const port = "8080"

type Data struct {
	Players []internal.Player
	Games   []internal.Game
}

func main() {
	fmt.Println("welcome to the program")

	e := echo.New()

	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	players := make([]internal.Player, 0)
	players = append(players, *internal.NewPlayer("jeremy"))

	data := Data{
		players,
		make([]internal.Game, 0),
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", data)
	})

	e.GET("/home", func(c echo.Context) error {
		return c.Render(200, "home", data)
	})

	e.GET("/player_stats", func(c echo.Context) error {
		return c.Render(200, "players", data)
	})

	e.GET("/games", func(c echo.Context) error {
		return c.Render(200, "games", data)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
