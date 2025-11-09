package model

type card interface {
	GetName() string
	GetType() cardType
}

type cardType string

const (
	PanicType cardType = "PANIC"
	ItemType  cardType = "ITEM"
)

type genericCard struct {
	Name string
	Type cardType
}

func (c genericCard) GetName() string   { return c.Name }
func (c genericCard) GetType() cardType { return c.Type }

func toCardSlice[T card](input []T) []card {
	res := make([]card, len(input))
	for i, v := range input {
		res[i] = v
	}
	return res
}
