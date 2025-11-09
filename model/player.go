package model

type player struct {
	Id             string
	OxygenCards    []card
	DiscardedCards []card
	Inventory      []card
	DiveLevel      int
}

func NewPlayer(id string) player {
	return player{
		Id: id,
	}
}

func (p player) Draw(cardType CardType, numberOfCards int) []genericCard {
	return make([]genericCard, 0)
}

func (p player) DoAction(state state) action {
	return NewAction(Dive, map[ActionParam]int{DiveLevels: 1})
}

func (p player) IsDead() bool {
	return len(p.OxygenCards) == 0
}
