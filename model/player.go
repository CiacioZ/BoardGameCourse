package model

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type playerEffect string

const (
	SkipTurn       playerEffect = "SKIP_TURN"
	CantMove       playerEffect = "CANT_MOVE"
	CantExplore    playerEffect = "CANT_EXPLORE"
	HaveToCalmDown playerEffect = "HAVE_TO_CALM_DOWN"
)

type player struct {
	Id               string
	OxygenCards      []card
	HandCards        []card
	DiscardedCards   []card
	Inventory        []*item
	DiscardedObjects []item
	DiveLevel        int
	ActiveEffects    map[playerEffect]int
}

func NewPlayer(id string, inventorySlot int) player {
	return player{
		Id:               id,
		OxygenCards:      make([]card, 0),
		HandCards:        make([]card, 0),
		DiscardedCards:   make([]card, 0),
		Inventory:        make([]*item, inventorySlot),
		DiscardedObjects: make([]item, 0),
		DiveLevel:        1,
		ActiveEffects:    make(map[playerEffect]int),
	}
}

func (p *player) Draw(numberOfCards int) []card {

	if len(p.OxygenCards) <= numberOfCards {
		numberOfCards = len(p.OxygenCards)
	}

	drawedCards := make([]card, numberOfCards)
	for i := 0; i < numberOfCards; i++ {
		drawedCard := p.OxygenCards[0]
		drawedCards[i] = drawedCard
		p.OxygenCards = p.OxygenCards[1:]
	}

	return drawedCards
}

func (p *player) Discard(cards []card) {
	p.DiscardedCards = append(p.DiscardedCards, cards...)
}

func (p *player) Breath() []card {

	var oxygenNeeded int
	switch {
	case p.DiveLevel >= 1 && p.DiveLevel <= 3:
		oxygenNeeded = 1
	case p.DiveLevel >= 4 && p.DiveLevel <= 6:
		oxygenNeeded = 2
	case p.DiveLevel >= 7 && p.DiveLevel <= 9:
		oxygenNeeded = 3
	}

	return p.Draw(oxygenNeeded)
}

func (p player) DecideActionToDo(gameState state, availableActions []actionType) action {

	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Player '%s'choose action: A=ascend, D=dive, E=explore, C=calm, U=use object H=Hold\n", p.Id)
		fmt.Printf("Answer:")
		reader.Scan()
		answer := strings.TrimSpace(strings.ToUpper(reader.Text()))

		readActionParam := strings.Split(answer, "-")

		switch strings.ToUpper(readActionParam[0]) {
		case "A":
			value, err := strconv.Atoi(readActionParam[1])
			if err != nil {
				fmt.Println("Action argument must be an integer")
			}
			if value < 1 || value > 3 {
				fmt.Println("Action argument must be between 1 and 3")
			}
			return NewAction(Ascend, map[actionParam]int{AscendLevels: value})
		case "D":
			value, err := strconv.Atoi(readActionParam[1])
			if err != nil {
				fmt.Println("Action argument must be an integer")
			}
			if value < 1 || value > 3 {
				fmt.Println("Action argument must be between 1 and 3")
			}
			return NewAction(Dive, map[actionParam]int{DiveLevels: value})
		case "E":
			value, err := strconv.Atoi(readActionParam[1])
			if err != nil {
				fmt.Println("Action argument must be an integer")
			}
			if value < 1 || value > 3 {
				fmt.Println("Action argument must be between 1 and 3")
			}
			return NewAction(Explore, map[actionParam]int{ExploreTime: value})
		case "C":
			return NewAction(CalmDown, map[actionParam]int{})
		case "U":
			value, err := strconv.Atoi(readActionParam[1])
			if err != nil {
				fmt.Println("Action argument must be an integer")
			}
			if value < 1 || value > len(p.Inventory) {
				fmt.Println("Action argument must be between 1 and 3")
			}
			return NewAction(UseObject, map[actionParam]int{ItemToUse: value})
		case "H":
			return NewAction(UseObject, map[actionParam]int{})
		default:
			fmt.Println(fmt.Printf("'%s' is not a valid action", readActionParam[0]))
		}
	}

}

func (p player) IsDead() bool {
	return len(p.OxygenCards) == 0
}

func (p player) CheckPanic() []panicEffect {
	panicEffects := make([]panicEffect, 0)

	panics := map[panicType]int{
		Blue:   0,
		Black:  0,
		Red:    0,
		Green:  0,
		Yellow: 0,
		Purple: 0,
	}

	for _, card := range p.HandCards {
		if card.GetType() == ItemType {
			continue
		}

		panicCard := card.(panicCard)
		for _, panicType := range panicCard.panicTypes {
			panics[panicType] = panics[panicType] + 1
		}
	}

	for panicType, counter := range panics {
		if counter >= 3 {
			panicEffects = append(panicEffects, panicActivationEffects[panicType][p.DiveLevel]...)
		}
	}

	return panicEffects
}

func (p *player) CheckPlayerEffects() []actionType {

	availableAction := []actionType{
		Ascend,
		Dive,
		Explore,
		CalmDown,
		UseObject,
		Distract,
	}

	notAvailableMoves := make([]actionType, 0)

	for effectType, value := range p.ActiveEffects {
		switch effectType {
		case SkipTurn:
			newValue := value - 1
			if newValue > 0 {
				p.ActiveEffects[effectType] = newValue
			} else {
				delete(p.ActiveEffects, effectType)
			}

			return make([]actionType, 0)
		case CantMove:
			notAvailableMoves = append(notAvailableMoves, Ascend, Dive)
			delete(p.ActiveEffects, effectType)
		case CantExplore:
			notAvailableMoves = append(notAvailableMoves, Ascend, Dive)
			delete(p.ActiveEffects, effectType)
		case HaveToCalmDown:
			delete(p.ActiveEffects, effectType)
			return []actionType{CalmDown}
		}
	}

	return SubtractSlices(availableAction, notAvailableMoves)
}

func SubtractSlices(fullActionTypeList, prohibitedActionTypes []actionType) []actionType {

	toRemove := make(map[actionType]bool)
	for _, v := range prohibitedActionTypes {
		toRemove[v] = true
	}

	var result []actionType
	for _, v := range fullActionTypeList {
		if !toRemove[v] {
			result = append(result, v)
		}
	}
	return result
}
