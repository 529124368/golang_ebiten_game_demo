package role

type Monster struct {
	X         float64
	Y         float64
	State     int
	Direction int
}

func NewMonster(x, y float64, state, dir int) *Monster {
	m := &Monster{
		X:         x,
		Y:         y,
		State:     state,
		Direction: dir,
	}
	return m
}

//TODO
func (m *Monster) DeadEvent() {

}

//TODO
func (m *Monster) Attack() {

}
