package internal

type Player struct {
	Id           int
	Name         string
	Games_played int
	Wins         int
	Losses       int
	Tricks       int
	Lone_hands   int
}

func NewPlayer(name string, id int) *Player {
	return &Player{
		id,
		name,
		0,
		0,
		0,
		0,
		0,
	}
}
