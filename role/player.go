package role

type Player struct {
	X         float64
	Y         float64
	State     int
	Direction int
	MouseX    int
	MouseY    int
}

func NewPlayer(x, y float64, state, dir, mx, my int) *Player {
	play := &Player{
		X:         x,
		Y:         y,
		State:     state,
		Direction: dir,
		MouseX:    mx,
		MouseY:    my,
	}
	return play
}

//TODO
func (p *Player) Attack() {

}

//TODO
func (p *Player) DeadEvent() {

}
