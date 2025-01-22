package internal

type Game struct {
	Id       int
	T1       [2]int
	T2       [2]int
	T1_names [2]string
	T2_names [2]string

	Team_points [2]int

	Winner     bool
	Win_by_two bool

	Rounds int
}
