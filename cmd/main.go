package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"

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

func sortPlayers(players []internal.Player) []internal.Player {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Tricks < players[j].Tricks
	})
	return players
}

func main() {
	player_id := 0

	fmt.Println("welcome to the program")

	e := echo.New()

	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	players := make([]internal.Player, 0)
	players = append(players, *internal.NewPlayer("jeremy", player_id))
	player_id++

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

	e.GET("/players/:id", func(c echo.Context) error {

		param := c.Param("id")

		id, err := strconv.Atoi(param)
		if err != nil {
			return c.Render(400, "players", data)
		}

		for _, val := range data.Players {
			if val.Id == id {
				return c.Render(200, "player_card", val)
			}
		}

		msg := fmt.Sprintf("No player with id: {%d}\n", id)
		return c.String(http.StatusNotFound, msg)
	})

	e.GET("/players", func(c echo.Context) error {

		json_data := echo.Map{}
		if err := c.Bind(&json_data); err != nil {
			return c.Render(400, "players", data)
		}

		id, err := strconv.Atoi(fmt.Sprintf("%v", json_data["p_id"]))
		if err != nil {
			return c.Render(400, "players", data)
		}

		for _, val := range data.Players {
			if val.Id == id {
				return c.Render(200, "player_card", val)
			}
		}

		msg := fmt.Sprintf("No player with id: {%d}\n", id)
		return c.String(http.StatusNotFound, msg)
	})

	e.GET("/players/edit", func(c echo.Context) error {

		json_data := echo.Map{}
		if err := c.Bind(&json_data); err != nil {
			return c.Render(400, "players", data)
		}

		p_id := fmt.Sprintf("%v", json_data["p_id"])

		id, err := strconv.Atoi(p_id)
		if err != nil {
			return c.Render(400, "players", data)
		}

		for _, val := range data.Players {
			if val.Id == id {
				return c.Render(200, "player_update_card", val)
			}
		}
		msg := fmt.Sprintf("No player with id: {%d}\n", id)
		return c.String(http.StatusNotFound, msg)
	})

	e.PUT("/players", func(c echo.Context) error {

		json_data := echo.Map{}
		if err := c.Bind(&json_data); err != nil {
			return c.Render(400, "players", data)
		}

		id, err := strconv.Atoi(fmt.Sprintf("%v", json_data["p_id"]))
		if err != nil {
			return c.Render(400, "players", data)
		}
		games_played, _ := strconv.Atoi(c.FormValue("games_played"))
		wins, _ := strconv.Atoi(c.FormValue("wins"))
		losses, _ := strconv.Atoi(c.FormValue("losses"))
		tricks, _ := strconv.Atoi(c.FormValue("tricks"))
		lone_h, _ := strconv.Atoi(c.FormValue("lone_hands"))

		var temp internal.Player
		var idx int
		for i, val := range data.Players {
			if val.Id == id {
				temp = val
				idx = i
			}
		}

		temp.Games_played = games_played
		temp.Wins = wins
		temp.Losses = losses
		temp.Lone_hands = lone_h
		temp.Tricks = tricks

		data.Players[idx] = temp

		return c.Render(200, "player_card", temp)
	})

	e.GET("/player_stats", func(c echo.Context) error {
		return c.Render(200, "players", data)
	})

	e.POST("/new_player", func(c echo.Context) error {
		name := c.FormValue("name")
		data.Players = append(data.Players, *internal.NewPlayer(name, player_id))
		player_id++
		return c.Render(200, "player_list", data)
	})

	e.GET("/games", func(c echo.Context) error {
		return c.Render(200, "games", data)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
