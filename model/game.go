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
		fmt.Printf("\t[LOG] Explore action: drawing %d cards\n", action.params[ExploreTime])
		cards := p.Draw(action.params[ExploreTime])
		fmt.Printf("\t[LOG] Drawn %d cards from oxygen deck\n", len(cards))
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
		fmt.Printf("\t[LOG] Found %d panic cards and %d items (at dive level %d)\n", len(panicCards), len(items), p.DiveLevel)
		for _, item := range items {
			fmt.Printf("\t[LOG] Processing item: %s (type: %s)\n", item.name, item.itemType)
			itemPlaced := false
			for i := 0; i < len(p.Inventory); i++ {
				slot := p.Inventory[i]
				if slot == nil {
					p.Inventory[i] = &item
					itemPlaced = true
					fmt.Printf("\t[LOG] Item '%s' placed in empty inventory slot %d\n", item.name, i)
					break // Break after placing item in empty slot
				} else {
					reader := bufio.NewScanner(os.Stdin)

					for {
						fmt.Printf("\tDo you want to keep item '%s' and drop item '%s'? (Y/N): ", item.name, slot.name)
						reader.Scan()
						answer := strings.TrimSpace(strings.ToUpper(reader.Text()))

						if answer == "Y" {
							fmt.Printf("\t[LOG] Replacing item '%s' with '%s' in slot %d\n", slot.name, item.name, i)
							p.Inventory[i] = &item
							itemPlaced = true
							break // Break from input loop and then from slot loop
						} else if answer == "N" {
							fmt.Printf("\t[LOG] Player declined to replace item '%s' in slot %d\n", slot.name, i)
							break // Break from input loop, continue to next slot
						} else {
							fmt.Printf("\tPlease answer with Y or N.\n")
							// Continue input loop to ask again
						}
					}
					if itemPlaced {
						break // Break from slot loop after placing item
					}
				}
			}
			if !itemPlaced {
				fmt.Printf("\t[LOG] Item '%s' was not placed (all slots full and player declined all replacements)\n", item.name)
			}
		}

		p.HandCards = append(p.HandCards, panicCards...)
		fmt.Printf("\t[LOG] Added %d panic cards to hand\n", len(panicCards))

	case Dive:
		oldLevel := p.DiveLevel
		p.DiveLevel += action.params[DiveLevels]
		fmt.Printf("\t[LOG] Dive action: level changed from %d to %d (dived %d levels)\n", oldLevel, p.DiveLevel, action.params[DiveLevels])
		cards := p.Draw(action.params[DiveLevels])
		fmt.Printf("\t[LOG] Drawn %d cards from oxygen deck\n", len(cards))
		panicCount := 0
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
				panicCount++
			}
		}
		fmt.Printf("\t[LOG] Added %d panic cards to hand\n", panicCount)
	case CalmDown:
		fmt.Printf("\t[LOG] CalmDown action: drawing 1 card\n")
		cards := p.Draw(1)
		panicCount := 0
		nonPanicCount := 0
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
				panicCount++
			} else {
				p.DiscardedCards = append(p.DiscardedCards, card)
				nonPanicCount++
			}
		}
		fmt.Printf("\t[LOG] Added %d panic cards to hand, discarded %d non-panic cards\n", panicCount, nonPanicCount)
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
						fmt.Printf("\t[LOG] Discarding panic card '%s' (%d/%d discarded)\n", panicCard.GetName(), discardedCard, 3)
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
		fmt.Printf("\t[LOG] CalmDown complete: discarded %d panic cards\n", discardedCard)
	case Ascend:
		oldLevel := p.DiveLevel
		p.DiveLevel -= action.params[AscendLevels]
		if p.DiveLevel < 1 {
			p.DiveLevel = 1
		}
		fmt.Printf("\t[LOG] Ascend action: level changed from %d to %d (ascended %d levels)\n", oldLevel, p.DiveLevel, action.params[AscendLevels])
		cards := p.Draw(action.params[AscendLevels])
		fmt.Printf("\t[LOG] Drawn %d cards from oxygen deck\n", len(cards))
		panicCount := 0
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
				panicCount++
			}
		}
		fmt.Printf("\t[LOG] Added %d panic cards to hand\n", panicCount)
		discardedCard := 0
		fmt.Printf("\t[LOG] Must discard %d panic cards\n", action.params[AscendLevels]+1)
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
						fmt.Printf("\t[LOG] Discarding panic card '%s' (%d/%d discarded)\n", panicCard.GetName(), discardedCard, action.params[AscendLevels]+1)
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
		fmt.Printf("\t[LOG] Ascend complete: discarded %d panic cards\n", discardedCard)
	case Distract:
		fmt.Printf("\t[LOG] Distract action: drawing 2 cards\n")
		cards := p.Draw(2)
		panicCount := 0
		for _, card := range cards {
			if card.GetType() == PanicType {
				p.HandCards = append(p.HandCards, card)
				panicCount++
			}
		}
		fmt.Printf("\t[LOG] Added %d panic cards to hand\n", panicCount)
		affectedPlayers := 0
		for _, player := range g.state.Players {
			if player.Id != p.Id && player.DiveLevel == p.DiveLevel {
				affectedPlayers++
				fmt.Printf("\t[LOG] Distracting player '%s' at same level (%d)\n", player.Id, player.DiveLevel)
				cards := player.Draw(2)
				playerPanicCount := 0
				for _, card := range cards {
					if card.GetType() == PanicType {
						player.HandCards = append(player.HandCards, card)
						playerPanicCount++
					}
				}
				fmt.Printf("\t[LOG] Player '%s' received %d panic cards\n", player.Id, playerPanicCount)
				value, found := player.ActiveEffects[CantExplore]
				if !found {
					player.ActiveEffects[CantExplore] = 1
					fmt.Printf("\t[LOG] Player '%s' now has CantExplore effect (1 turn)\n", player.Id)
				} else {
					player.ActiveEffects[CantExplore] = value + 1
					fmt.Printf("\t[LOG] Player '%s' CantExplore effect extended to %d turns\n", player.Id, value+1)
				}
			}
		}
		if affectedPlayers == 0 {
			fmt.Printf("\t[LOG] No other players at same level to distract\n")
		}

	case UseObject:
		itemToUse, hasItemParam := action.params[ItemToUse]
		if !hasItemParam {
			fmt.Printf("\t[LOG] UseObject action: Hold action (no item specified)\n")
			return
		}
		itemIndex := itemToUse - 1 // Convert to 0-based index
		if itemIndex < 0 || itemIndex >= len(p.Inventory) {
			fmt.Printf("\t[LOG] UseObject action: invalid item index %d\n", itemToUse)
			return
		}
		itemToActivate := p.Inventory[itemIndex]
		if itemToActivate == nil {
			fmt.Printf("\t[LOG] UseObject action: no item at inventory slot %d\n", itemToUse)
			return
		}
		fmt.Printf("\t[LOG] UseObject action: using item '%s' (slot %d)\n", itemToActivate.name, itemToUse)
		effectCount := 0
		for _, effect := range itemToActivate.effects {
			effectCount++
			fmt.Printf("\t[LOG] Applying effect: %s (value: %d)\n", effect.effectType, effect.value)
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
		if effectCount == 0 {
			fmt.Printf("\t[LOG] Item '%s' has no effects\n", itemToActivate.name)
		}
	}

}

func printCards(cards []card, message string) {
	fmt.Println(message)
	for _, card := range cards {
		fmt.Printf("\t\t%s\n", card.GetName())
	}
}
