package internal

type Game struct {
	Players      []Player
	Teams        []Player
	Teams_points []int

	Winner     bool
	Win_by_two bool

	Rounds int
}
