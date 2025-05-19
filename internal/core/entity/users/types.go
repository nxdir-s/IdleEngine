package users

type State struct {
	Exp uint32
}

func NewState() *State {
	return &State{}
}
