package internal

type Game struct {
	players      []Player
	teams        []Player
	teams_points []int

	winner     bool
	win_by_two bool

	rounds int
}
