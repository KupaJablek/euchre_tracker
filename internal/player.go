package internal

type Player struct {
	Name         string
	Games_playes int
	Wins         int
	Losses       int
	Tricks       int
    Lone_hands   int
    Points       int
}

func NewPlayer(name string) *Player {
	return &Player{
		name,
		0,
		0,
		0,
		0,
		0,
        0,
	}
}
