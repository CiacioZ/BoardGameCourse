package model

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type GameParameterType string

const (
	NumberOfPlayers                    GameParameterType = "NUMBER_OF_PLAYERS"
	NumberOfPanicCardsToActivateEffect GameParameterType = "NUMBER_OF_PANIC_CARD_TO_ACTIVATE_EFFECT"
	NumberOfItemSlots                  GameParameterType = "NUMBER_OF_ITEM_SLOTS"
	NumberOfAmuletsToWin               GameParameterType = "NUMBER_OF_AMULETS_TO_WIN"
)

type gameParameter struct {
	parameterType GameParameterType
	value         int
}

func NewGameParameter(parameterType GameParameterType, value int) gameParameter {
	return gameParameter{
		parameterType: parameterType,
		value:         value,
	}
}

type game struct {
	state      state
	parameters parameters
	randomizer *rand.Rand
}

type parameters struct {
	values map[GameParameterType]int
}

func NewGameParameters(params []gameParameter) parameters {
	config := parameters{
		values: make(map[GameParameterType]int),
	}

	for _, param := range params {
		config.values[param.parameterType] = param.value
	}

	return config
}

type state struct {
	Players      []player
	Round        int
	ActualPlayer int
}

func NewState() state {
	return state{
		Players:      make([]player, 0),
		Round:        0,
		ActualPlayer: -1,
	}
}

func (s *state) AddPlayer(player player) {
	s.Players = append(s.Players, player)
}

func (s *state) SetRound(i int) {
	s.Round = i
}

func (s *state) NextRound() {
	s.Round++
}

func (s *state) NextPlayer() {
	s.ActualPlayer += 1
	if s.ActualPlayer == len(s.Players) {
		s.ActualPlayer = 0
		s.NextRound()
	}
}

func NewGame(
	parameters ...gameParameter,
) game {

	g := game{
		parameters: NewGameParameters(parameters),
		state:      NewState(),
		randomizer: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for i := range g.parameters.values[NumberOfPlayers] {
		player := NewPlayer(fmt.Sprintf("P%d", i+1), g.parameters.values[NumberOfItemSlots])

		player.OxygenCards = g.GenerateOxygenDeck()

		g.state.AddPlayer(player)
	}

	return g
}

func (g *game) Run(numberOfGames int) {

	g.state.Round = 1

	for !g.IsGameEnded() {
		fmt.Printf("Start round: %d\n", g.state.Round)

		for i := 0; i < len(g.state.Players); i++ {
			g.state.ActualPlayer = i

			p := g.GetActualPlayer()

			fmt.Printf("Player '%s':\n", p.Id)

			fmt.Printf("\tStart turn\n")
			fmt.Printf("\tLevel %d\n", p.DiveLevel)
			fmt.Printf("\tInventory:\n")
			for _, item := range p.Inventory {
				if item != nil {
					fmt.Printf("\t\t%s\n", item.name)
				}
			}

			//BREATH
			cards := p.Breath()
			printCards(cards, "\tBreath:\n")
			for _, card := range cards {
				if card.GetType() == PanicType {
					p.HandCards = append(p.HandCards, card)
				} else {
					p.DiscardedCards = append(p.DiscardedCards, card)
				}
			}
			if p.IsDead() {
				g.state.NextPlayer()
				continue
			}

			printCards(p.HandCards, "\tHand:\n")

			//CHECK PANIC
			effects := p.CheckPanic()
			if len(effects) > 0 {
				fmt.Printf("\tApply Effects:\n")
				for _, effect := range effects {
					fmt.Printf("\t\t%s", effect.effectType)
				}
				g.ApplyEffect(&g.state.Players[g.state.ActualPlayer], effects)
			}

			//CHECK PLAYER EFFECTS
			availableActions := p.CheckPlayerEffects()
			fmt.Printf("\tAvailable Actions: %+v\n", availableActions)

			//DECIDE ACTION TO DO
			actionToDo := p.DecideActionToDo(g.state, availableActions)
			fmt.Printf("\tAction To Do: %+v\n", actionToDo)

			//RESOLVE ACTION
			g.resolveAction(p, actionToDo)
			fmt.Printf("\tAction resolved\n")

			fmt.Printf("\tInventory:\n")
			for _, item := range p.Inventory {
				if item != nil {
					fmt.Printf("\t\t%s\n", item.name)
				}
			}
			printCards(p.HandCards, "\tHand:\n")

			//CHECK PANIC
			effects = p.CheckPanic()
			if len(effects) > 0 {
				fmt.Printf("\tApply Effects:\n")
				for _, effect := range effects {
					fmt.Printf("\t\t%s\n", effect.effectType)
				}
				g.ApplyEffect(&g.state.Players[g.state.ActualPlayer], effects)
			}

			fmt.Printf("\tEnd turn\n\n")
		}

		fmt.Printf("End round: %d\n", g.state.Round)

		g.state.NextRound()
	}
}

func (g game) GetActualPlayer() *player {
	return &g.state.Players[g.state.ActualPlayer]
}

func (g *game) IsDavyJonesIsDead() bool {
	amuletAtLevel10 := 0

	for _, player := range g.state.Players {
		if player.DiveLevel == 10 {
			for _, item := range player.Inventory {
				if item.itemType == Amulets {
					amuletAtLevel10++
				}
			}
		}
	}

	return amuletAtLevel10 >= 3
}

func (g *game) AreAllPlayersDead() bool {

	AllDead := true
	for _, player := range g.state.Players {
		if !player.IsDead() {
			AllDead = false
			break
		}
	}

	return AllDead
}

func (g *game) IsGameEnded() bool {
	return g.AreAllPlayersDead() || g.IsDavyJonesIsDead()
}

func (g *game) GenerateOxygenDeck() []card {

	deck := make([]card, 0)

	//Generate one color type panic card
	for i := 1; i <= 2; i++ {
		deck = append(deck, toCardSlice(singlePanicCards)...)

	}

	//Generate two color type panic card
	deck = append(deck, toCardSlice(doublePanicCards)...)

	//Generate three color type panic card
	deck = append(deck, toCardSlice(triplePanicCards)...)

	//Generate common item card
	deck = append(deck, toCardSlice(commonItemCards)...)

	//Generate uncommon item card
	deck = append(deck, toCardSlice(uncommonItemCards)...)

	//Generate rare item card
	deck = append(deck, toCardSlice(rareItemCards)...)

	//Generate legendary item card
	deck = append(deck, toCardSlice(legendaryItemCards)...)

	g.randomizer.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

func (g *game) ApplyEffect(p *player, effects []panicEffect) {
	for _, effect := range effects {
		switch effect.effectType {
		case MoveUp:
			p.DiveLevel -= effect.value
			if p.DiveLevel < 1 {
				p.DiveLevel = 1
			}
		case MoveDown:
			p.DiveLevel += effect.value
			if p.DiveLevel > 10 {
				p.DiveLevel = 10
			}
		case CannotExplore:
			p.ActiveEffects[CantExplore] = effect.value
		case JumpTurn:
			p.ActiveEffects[SkipTurn] = effect.value
		case DiscardO2:
			cards := p.Draw(effect.value)
			p.Discard(cards)
		case DropObject:
			for i := 0; i < effect.value && i < len(p.Inventory); i++ {
				item := p.Inventory[i]
				if item != nil && item.itemType == Utility {
					p.DiscardedObjects = append(p.DiscardedObjects, *item)
					p.Inventory[i] = nil
				}
			}
		case MustCalmDown:
			p.ActiveEffects[HaveToCalmDown] = effect.value
		case MoveToFreeLevel:
			if p.DiveLevel > 1 {
				for destinationLevel := p.DiveLevel - 1; destinationLevel <= 1; destinationLevel-- {
					for _, player := range g.state.Players {
						if player.Id != p.Id && player.DiveLevel == destinationLevel {
							continue
						}
					}
				}
			}

		case DropTreasureToken:
			dropped := 0
			for i := 0; i < len(p.Inventory) && dropped < effect.value; i++ {
				item := p.Inventory[i]
				if item != nil && item.itemType == TreasureToken {
					p.DiscardedObjects = append(p.DiscardedObjects, *item)
					p.Inventory[i] = nil
					dropped++
				}
			}
		case DropAmulet:
			for i := 0; i < effect.value && i < len(p.Inventory); i++ {
				item := p.Inventory[i]
				if item != nil && item.itemType == Amulets {
					p.DiscardedObjects = append(p.DiscardedObjects, *item)
					p.Inventory[i] = nil
				}
			}
		case DropEverything:
			for i := 0; i < effect.value && i < len(p.Inventory); i++ {
				item := p.Inventory[i]
				if item != nil {
					p.DiscardedObjects = append(p.DiscardedObjects, *item)
					p.Inventory[i] = nil
				}
			}
		case DropEverythingButAmulets:
			for i := 0; i < effect.value && i < len(p.Inventory); i++ {
				item := p.Inventory[i]
				if item != nil && item.itemType != Amulets {
					p.DiscardedObjects = append(p.DiscardedObjects, *item)
					p.Inventory[i] = nil
				}
			}
		case DropO2ForSameLevelPlayers:
			for _, player := range g.state.Players {
				if player.Id != p.Id && player.DiveLevel == p.DiveLevel {
					cards := player.Draw(effect.value)
					player.HandCards = append(player.HandCards, cards...)
				}
			}
		case DrawO2:
			cards := p.Draw(effect.value)
			p.HandCards = append(p.HandCards, cards...)
		}
	}
}

func (g *game) resolveAction(p *player, action action) {

	switch action.actionType {

	case Explore:
		cards := p.Draw(action.params[ExploreTime])
		panicCards := make([]card, 0)
		items := make([]item, 0)
		for _, card := range cards {
			if card.GetType() == PanicType {
				panicCards = append(panicCards, card)
			} else {
				itemCard := card.(itemCard)
				items = append(items, itemCard.items[p.DiveLevel])
			}
		}
		for _, item := range items {
			for i, slot := range p.Inventory {
				if slot == nil {
					p.Inventory[i] = &item
				} else {
					reader := bufio.NewScanner(os.Stdin)

					fmt.Printf("\tDo you want to keep item '%s' and drop item '%s'? (Y/N): ", item.name, slot.name)
					reader.Scan()
					answer := strings.TrimSpace(strings.ToUpper(reader.Text()))

					if answer == "Y" {
						p.Inventory[i] = &item
						break
					} else if answer == "N" {
						continue
					} else {
						fmt.Printf("\tPlease answer with Y or N.\n")
					}

				}
			}
		}

		p.HandCards = append(p.HandCards, panicCards...)

	case Dive:
		p.DiveLevel += action.params[DiveLevels]
		cards := p.Draw(action.params[DiveLevels])
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
			}
		}
	case CalmDown:
		cards := p.Draw(1)
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
			} else {
				p.DiscardedCards = append(p.DiscardedCards, card)
			}
		}
		discardedCard := 0
		for discardedCard < 3 && len(p.HandCards) > 0 {
			removed := false
			// Iterate backwards to safely remove elements
			for i := len(p.HandCards) - 1; i >= 0; i-- {
				panicCard := p.HandCards[i]
				reader := bufio.NewScanner(os.Stdin)

				for {
					fmt.Printf("\tDo you want to discard panic card '%s'? (Y/N): ", panicCard.GetName())
					reader.Scan()
					answer := strings.TrimSpace(strings.ToUpper(reader.Text()))

					if answer == "Y" {
						discardedCard++
						p.DiscardedCards = append(p.DiscardedCards, panicCard)
						p.HandCards = append(p.HandCards[:i], p.HandCards[i+1:]...)
						removed = true
						break // Break from input loop
					} else if answer == "N" {
						break // Break from input loop to continue to next card
					} else {
						fmt.Printf("\tPlease answer with Y or N.\n")
						// Continue loop to ask again
					}
				}
				// If we removed a card, break to restart the outer loop
				if removed {
					break
				}
			}
		}
	case Ascend:
		p.DiveLevel -= action.params[AscendLevels]
		cards := p.Draw(action.params[AscendLevels])
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
			}
		}
		discardedCard := 0
		for discardedCard < action.params[AscendLevels]+1 && len(p.HandCards) > 0 {
			removed := false
			// Iterate backwards to safely remove elements
			for i := len(p.HandCards) - 1; i >= 0; i-- {
				panicCard := p.HandCards[i]
				reader := bufio.NewScanner(os.Stdin)

				for {
					fmt.Printf("\tDo you want to discard panic card '%s'? (Y/N): ", panicCard.GetName())
					reader.Scan()
					answer := strings.TrimSpace(strings.ToUpper(reader.Text()))

					if answer == "Y" {
						discardedCard++
						p.DiscardedCards = append(p.DiscardedCards, panicCard)
						p.HandCards = append(p.HandCards[:i], p.HandCards[i+1:]...)
						removed = true
						break // Break from input loop
					} else if answer == "N" {
						break // Break from input loop to continue to next card
					} else {
						fmt.Printf("\tPlease answer with Y or N.\n")
						// Continue loop to ask again
					}
				}
				// If we removed a card, break to restart the outer loop
				if removed {
					break
				}
			}
		}
	case Distract:
		cards := p.Draw(2)
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
			}
		}
		for _, player := range g.state.Players {
			if player.Id != p.Id && player.DiveLevel == p.DiveLevel {
				cards := player.Draw(2)
				for _, card := range cards {
					if card.GetType() == PanicType {
						player.HandCards = append(player.HandCards, card)
					}
				}
				value, found := player.ActiveEffects[CantExplore]
				if !found {
					player.ActiveEffects[CantExplore] = 1
				} else {
					player.ActiveEffects[CantExplore] = value + 1
				}
			}
		}

	case UseObject:
		itemToActivate := p.Inventory[action.params[ItemToUse]]
		for _, effect := range itemToActivate.effects {
			switch effect.effectType {
			case LookNextO2Cards:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case MovementCostReduction:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case BreathCostReduction:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case BlockPlayer:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case IgnorePanicActivation:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case AnotherPlayerMustDrawO2:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case StealItemFromPlayer:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case StealAmuletFromPLayer:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case RecoverDiscardedO2:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			case ReorderNextO2Cards:
				fmt.Printf("%s NOT IMPLEMENTED", effect.effectType)
			}
		}
	}

}

func printCards(cards []card, message string) {
	fmt.Println(message)
	for _, card := range cards {
		fmt.Printf("\t\t%s\n", card.GetName())
	}
}
