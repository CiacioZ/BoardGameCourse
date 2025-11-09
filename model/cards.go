package model

type card interface {
	GetName() string
	GetType() CardType
}

type CardType string

const (
	PanicType CardType = "PANIC"
	ItemType  CardType = "ITEM"
)

type genericCard struct {
	Name string
	Type CardType
}

func (c genericCard) GetName() string   { return c.Name }
func (c genericCard) GetType() CardType { return c.Type }
