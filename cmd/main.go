package main

import (
	"errors"
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

var data Data

func getPlayer(id int) (internal.Player, error) {
	for _, player := range data.Players {
		if player.Id == id {
			return player, nil
		}
	}

	return internal.Player{}, errors.New("player id not found")
}

func getName(id int) string {
	for _, player := range data.Players {
		if player.Id == id {
			return player.Name
		}
	}

	return ""
}

func main() {
	player_id := 0
	game_id := 0

	fmt.Println("welcome to the program")

	e := echo.New()

	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	testers := []string{
		"Lisa", "Billy", "Trish", "Mike",
		"Kate", "Ed", "Melissa", "Grandpa",
		"Gramdma", "Evan", "Tracy", "Alvaro",
		"Maximus", "Tamara", "Jane",
	}

	players := make([]internal.Player, 0)
	for _, person := range testers {
		players = append(players, *internal.NewPlayer(person, player_id))
		player_id++
	}

	data = Data{
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

	e.GET("/players/edit/:id", func(c echo.Context) error {

		p_id := c.Param("id")

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

	e.PUT("/players/:id", func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
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

	e.POST("/games", func(c echo.Context) error {
		t1_p1, _ := strconv.Atoi(c.FormValue("t1_p1"))
		t1_p2, _ := strconv.Atoi(c.FormValue("t1_p2"))
		t2_p1, _ := strconv.Atoi(c.FormValue("t2_p1"))
		t2_p2, _ := strconv.Atoi(c.FormValue("t2_p2"))

		// validate game before creation

		// setup new game
		temp := internal.Game{}

		temp.Id = game_id
		game_id++
		temp.T1 = [2]int{t1_p1, t1_p2}
		temp.T2 = [2]int{t2_p1, t2_p2}

		temp.T1_names[0] = getName(t1_p1)
		temp.T1_names[1] = getName(t1_p2)
		temp.T2_names[0] = getName(t2_p1)
		temp.T2_names[1] = getName(t2_p2)

		temp.Team_points = [2]int{0, 0}

		data.Games = append(data.Games, temp)

		return c.Render(200, "games_list", data)
	})

	e.GET("/games", func(c echo.Context) error {
		return c.Render(200, "games", data)
	})

	e.GET("/play_game/:id", func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.Render(200, "games_page", data)
		}

		for _, game := range data.Games {
			if game.Id == id {
				return c.Render(200, "play_game", game)
			}
		}

		return c.Render(200, "games_page", data)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
