package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// Color helper functions
func colorize(text string, color string) string {
	return color + text + colorReset
}

func red(text string) string {
	return colorize(text, colorRed)
}

func green(text string) string {
	return colorize(text, colorGreen)
}

func yellow(text string) string {
	return colorize(text, colorYellow)
}

func blue(text string) string {
	return colorize(text, colorBlue)
}

func cyan(text string) string {
	return colorize(text, colorCyan)
}

func bold(text string) string {
	return colorize(text, colorBold)
}

var NumberOfGames = 1
var NumberOfPlayers = 2

type abilityType string

const (
	Encounter     abilityType = "ENCOUNTER"
	Environment   abilityType = "ENVIRONMENT"
	Technical     abilityType = "TECHNICAL"
	Soprannatural abilityType = "SOPRANNATURAL"
)

type O2 struct {
	Type       abilityType
	Value      int
	ItemReward int
}

type player struct {
	Id               string
	O2               []O2
	DiscardedO2      []O2
	Ability          map[abilityType]abilityValue
	Panic            map[abilityType]int
	PanicTollerance  map[abilityType]int
	Treasure         int
	Inventory        []item
	Effects          []itemEffect
	MaxInventorySize int
}

type abilityValue struct {
	Value    int
	Modifier int
}

type item struct {
	Type    string
	Effects []itemEffect
}

type effectType string

const (
	IncreaseEncounter     effectType = "INCREASE_ENCOUNTER"
	IncreaseEnvironment   effectType = "INCREASE_ENVIRONMENT"
	IncreaseTechnical     effectType = "INCREASE_TECHNICAL"
	IncreaseSoprannatural effectType = "INCREASE_SOPRANNATURAL"
	DrawMoreItems         effectType = "DRAW_MORE_ITEMS"
	Treasure              effectType = "TREASURE"
	AddSlotInventory      effectType = "ADD_SLOT_INVENTORY"
	RecoverDiscardedCards effectType = "RECOVER_DISCARDED_CARDS"
	FreeBreath            effectType = "FREE_BREATH"
	ReduceOnePanicType    effectType = "REDUCE_ONE_PANIC"
	LookAndReorder        effectType = "LOOK_AND_REORDER"
	ReduceAllPanicTypes   effectType = "REDUCE_ALL_PANIC"
)

type itemEffect struct {
	effectType effectType
	value      int
}

// hasFreeBreath checks if the player has the FreeBreath effect active
func hasFreeBreath(p *player) bool {
	for _, effect := range p.Effects {
		if effect.effectType == FreeBreath {
			return true
		}
	}
	return false
}

// calculateItemsToDraw calculates how many items to draw including DrawMoreItems effects
func calculateItemsToDraw(baseReward int, p *player) int {
	itemsToDraw := baseReward
	for _, effect := range p.Effects {
		if effect.effectType == DrawMoreItems {
			itemsToDraw += effect.value
		}
	}
	return itemsToDraw
}

// reducePanic reduces panic by 1 and ensures it doesn't go below 0
func reducePanic(p *player, panicType abilityType) {
	if p.Panic[panicType] > 0 {
		p.Panic[panicType] -= 1
		if p.Panic[panicType] < 0 {
			p.Panic[panicType] = 0
		}
	}
}

// reducePanicBy reduces panic by a specific value and ensures it doesn't go below 0
func reducePanicBy(p *player, panicType abilityType, value int) {
	p.Panic[panicType] = p.Panic[panicType] - value
	if p.Panic[panicType] < 0 {
		p.Panic[panicType] = 0
	}
}

// handlePanicTrigger handles all effects when panic is triggered
func handlePanicTrigger(p *player, panicType abilityType) {
	fmt.Printf("\t\t\t%s\n", red(bold("*** PANIC! ***")))
	p.Panic[panicType] = 0

	// Lose all items
	fmt.Printf("\t\t\t%s\n", yellow("Lost all items from inventory"))
	p.Inventory = make([]item, 0)

	// Discard 5 O2 cards
	cardsToDiscard, newO2AfterPanic := draw(5, p.O2)
	p.O2 = newO2AfterPanic
	p.DiscardedO2 = append(p.DiscardedO2, cardsToDiscard...)
	fmt.Printf("\t\t\t%s\n", yellow(fmt.Sprintf("Discarded %d O2 cards due to panic", len(cardsToDiscard))))
}

// drawO2CardForBreath draws an O2 card for breathing, handling FreeBreath effect
func drawO2CardForBreath(p *player) ([]O2, bool) {
	var cards []O2
	var newO2 []O2

	if !hasFreeBreath(p) {
		cards, newO2 = draw(1, p.O2)
		if len(cards) == 0 {
			return nil, false
		}
		p.O2 = newO2
		p.DiscardedO2 = append(p.DiscardedO2, cards...)
	} else {
		// FreeBreath: draw card but put it back in the deck instead of discarding
		cards, newO2 = draw(1, p.O2)
		if len(cards) == 0 {
			return nil, false
		}
		p.O2 = newO2
		// Card is put back at the bottom of the deck with FreeBreath
		p.O2 = append(p.O2, cards[0])
	}
	return cards, true
}

// checkCardResolved checks if a card is resolved successfully
func checkCardResolved(p *player, card O2, diceResult int) bool {
	ability := p.Ability[card.Type]
	return ability.Value+ability.Modifier+diceResult > card.Value
}

// increasePanic increases panic by 1 and checks if panic is triggered
func increasePanic(p *player, panicType abilityType) bool {
	p.Panic[panicType] += 1
	return p.Panic[panicType] >= p.PanicTollerance[panicType]
}

// addItemsToInventory adds items to inventory, handling full inventory and treasure items
func addItemsToInventory(p *player, items []item, actionType string, cardType abilityType, abilityValue, modifier, diceResult, cardValue int) {
	for _, item := range items {
		// Check if item is a treasure (has Treasure effect)
		isTreasure := false
		treasureValue := 0
		for _, effect := range item.Effects {
			if effect.effectType == Treasure {
				isTreasure = true
				treasureValue += effect.value
			}
		}

		if isTreasure {
			// Treasure items go directly to treasure, not inventory
			p.Treasure += treasureValue
			fmt.Printf("\t\t\t%s '%s' %s %d + %d + %d > %d: found treasure = %s %s\n",
				actionType, cardType, green("RESOLVED"), abilityValue, modifier, diceResult, cardValue,
				yellow(item.Type), green(fmt.Sprintf("(+%d treasure)", treasureValue)))
		} else {
			// Regular items go to inventory
			if len(p.Inventory) < p.MaxInventorySize {
				p.Inventory = append(p.Inventory, item)
				fmt.Printf("\t\t\t%s '%s' %s %d + %d + %d > %d: item found = %s\n",
					actionType, cardType, green("RESOLVED"), abilityValue, modifier, diceResult, cardValue,
					cyan(item.Type))
			} else {
				fmt.Printf("\t\t\t%s '%s' %s %d + %d + %d > %d: item found = %s %s\n",
					actionType, cardType, green("RESOLVED"), abilityValue, modifier, diceResult, cardValue,
					cyan(item.Type), yellow("(INVENTORY FULL, ITEM LOST)"))
			}
		}
	}
}

// increaseAbilityModifier increases the modifier for a specific ability type
func increaseAbilityModifier(p *player, abilityType abilityType, value int) {
	ability := p.Ability[abilityType]
	ability.Modifier += value
	p.Ability[abilityType] = ability
}

// resetTemporaryEffects resets all temporary effects at the end of a turn
func resetTemporaryEffects(p *player) {
	// Reset ability modifiers to 0
	for abilityType := range p.Ability {
		ability := p.Ability[abilityType]
		ability.Modifier = 0
		p.Ability[abilityType] = ability
	}

	// Clear temporary effects (DrawMoreItems, FreeBreath)
	// Keep only permanent effects if any (currently none are permanent in Effects)
	p.Effects = make([]itemEffect, 0)
}

// printPlayerStatus prints the current status of the player
func printPlayerStatus(p *player) {
	fmt.Printf("\t\t\t%s\n", bold(cyan(fmt.Sprintf("--- %s STATUS ---", p.Id))))
	fmt.Printf("\t\t\tO2 Cards: %s | Discarded: %s\n",
		cyan(fmt.Sprintf("%d", len(p.O2))), yellow(fmt.Sprintf("%d", len(p.DiscardedO2))))
	fmt.Printf("\t\t\tTreasure: %s\n", yellow(fmt.Sprintf("%d", p.Treasure)))
	fmt.Printf("\t\t\tInventory (%s/%s): ",
		cyan(fmt.Sprintf("%d", len(p.Inventory))), cyan(fmt.Sprintf("%d", p.MaxInventorySize)))
	if len(p.Inventory) == 0 {
		fmt.Print(yellow("Empty"))
	} else {
		for i, item := range p.Inventory {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(cyan(item.Type))
		}
	}
	fmt.Println()
	fmt.Printf("\t\t\tAbilities: ")
	first := true
	for abilityType, ability := range p.Ability {
		if !first {
			fmt.Print(" | ")
		}
		fmt.Printf("%s: %s", abilityType, cyan(fmt.Sprintf("%d", ability.Value)))
		if ability.Modifier != 0 {
			if ability.Modifier > 0 {
				fmt.Printf(" %s", green(fmt.Sprintf("(+%d)", ability.Modifier)))
			} else {
				fmt.Printf(" %s", red(fmt.Sprintf("(%d)", ability.Modifier)))
			}
		}
		first = false
	}
	fmt.Println()
	fmt.Printf("\t\t\tPanic: ")
	first = true
	for panicType, panicValue := range p.Panic {
		if !first {
			fmt.Print(" | ")
		}
		panicColor := green
		if panicValue >= p.PanicTollerance[panicType] {
			panicColor = red
		} else if panicValue > 0 {
			panicColor = yellow
		}
		fmt.Printf("%s: %s/%s", panicType,
			panicColor(fmt.Sprintf("%d", panicValue)),
			cyan(fmt.Sprintf("%d", p.PanicTollerance[panicType])))
		first = false
	}
	fmt.Println()
	fmt.Printf("\t\t\t%s\n", cyan("-------------------"))
}

func main() {

	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	commonItemsCount := 16
	uncommonItemsCount := 8
	rareItemsCount := 6
	veryRareItemsCount := 3
	legendaryItemsCount := 1

	commonItems := []item{
		{
			Type: "Coltello", // +1 agli incontri
			Effects: []itemEffect{
				{
					effectType: IncreaseEncounter,
					value:      2,
				},
			},
		},
		{
			Type: "Torcia", // +1 agli oggetti trovati
			Effects: []itemEffect{
				{
					effectType: DrawMoreItems,
					value:      1,
				},
			},
		},
		{
			Type: "Giubbotto rinforzato", // +1 agli incontri
			Effects: []itemEffect{
				{
					effectType: IncreaseEncounter,
					value:      1,
				},
			},
		},
		{
			Type: "1 moneta", // tesoro
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      1,
				},
			},
		},
		{
			Type: "Sacca maggiorata", // aggiunge uno slot in più all'inventario
			Effects: []itemEffect{
				{
					effectType: AddSlotInventory,
					value:      1,
				},
			},
		},
		{
			Type: "Attrezzi da riparazione", // +1 agli imprevisti tecnici
			Effects: []itemEffect{
				{
					effectType: IncreaseTechnical,
					value:      1,
				},
			},
		},
		{
			Type: "Pinne migliorate", // +1 agli imprevisti ambientali
			Effects: []itemEffect{
				{
					effectType: IncreaseEnvironment,
					value:      1,
				},
			},
		},
		{
			Type: "Bolla d'aria", // recupera una carta ossigeno scartata
			Effects: []itemEffect{
				{
					effectType: RecoverDiscardedCards,
					value:      1,
				},
			},
		},
		{
			Type: "Collana delle sirene", // +1 agli ineffetti soprannaturali
			Effects: []itemEffect{
				{
					effectType: IncreaseSoprannatural,
					value:      1,
				},
			},
		},
		{
			Type: "Torcia", // +1 agli oggetti trovati
			Effects: []itemEffect{
				{
					effectType: DrawMoreItems,
					value:      1,
				},
			},
		},
		{
			Type: "Amuleto del tritone", // +1 agli effetti soprannaturali
			Effects: []itemEffect{
				{
					effectType: IncreaseSoprannatural,
					value:      1,
				},
			},
		},
		{
			Type: "1 moneta", // tesoro
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      1,
				},
			},
		},
		{
			Type: "Sacca maggiorata", // aggiunge uno slot in più all'inventario
			Effects: []itemEffect{
				{
					effectType: AddSlotInventory,
					value:      1,
				},
			},
		},
		{
			Type: "Attrezzi da riparazione", // +1 alle esplorazioni
			Effects: []itemEffect{
				{
					effectType: IncreaseTechnical,
					value:      1,
				},
			},
		},
		{
			Type: "Pinne migliorate", // +1 agli imprevisti ambientali
			Effects: []itemEffect{
				{
					effectType: IncreaseEnvironment,
					value:      1,
				},
			},
		},
		{
			Type: "1 Moneta", // recupera una carta ossigeno scartata
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      1,
				},
			},
		},
	}

	uncommonItems := []item{
		{
			Type: "Fiocina", // +3 agli incontri
			Effects: []itemEffect{
				{
					effectType: IncreaseEncounter,
					value:      3,
				},
			},
		},
		{
			Type: "Barometro migliorato", // +3 agli imprevisti ambientali
			Effects: []itemEffect{
				{
					effectType: IncreaseEnvironment,
					value:      3,
				},
			},
		},
		{
			Type: "3 Monete", // tesoro
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      3,
				},
			},
		},
		{
			Type: "Collana di corallo", // +3 agli imprevisti ambientali
			Effects: []itemEffect{
				{
					effectType: IncreaseSoprannatural,
					value:      3,
				},
			},
		},
		{
			Type: "3 Monete", // tesoro
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      3,
				},
			},
		},
		{
			Type: "Sacca d'aria", // recupera 2 carte ossigeno scartate
			Effects: []itemEffect{
				{
					effectType: RecoverDiscardedCards,
					value:      2,
				},
			},
		},
		{
			Type: "Respiratore migliorato", // il respiro non consuma ossigeno
			Effects: []itemEffect{
				{
					effectType: FreeBreath,
					value:      1,
				},
			},
		},
		{
			Type: "Kit antistress", // riduce un tipo di panico di 1
			Effects: []itemEffect{
				{
					effectType: ReduceOnePanicType,
					value:      1,
				},
			},
		},
	}

	rareItems := []item{
		{
			Type: "Bombola aggiuntiva", // recupera 3 carte ossigeno scartate
			Effects: []itemEffect{
				{
					effectType: RecoverDiscardedCards,
					value:      3,
				},
			},
		},
		{
			Type: "Tridente di poseidone", // +2 agli incontri e agli imprevisti sovrannaturali
			Effects: []itemEffect{
				{
					effectType: IncreaseEncounter,
					value:      2,
				},
				{
					effectType: IncreaseSoprannatural,
					value:      2,
				},
			},
		},
		{
			Type: "Drone autoguidato", // +2 agli imprevisti ambientali e tecnici
			Effects: []itemEffect{
				{
					effectType: IncreaseEnvironment,
					value:      2,
				},
				{
					effectType: IncreaseTechnical,
					value:      2,
				},
			},
		},
		{
			Type: "Kit antistress migliorato", // riduce di 2 un tipo di panico
			Effects: []itemEffect{
				{
					effectType: ReduceOnePanicType,
					value:      2,
				},
			},
		},
		{
			Type: "5 Monete", // tesoro
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      5,
				},
			},
		},
		{
			Type: "5 Monete", // tesoro
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      5,
				},
			},
		},
	}

	veryRareItems := []item{
		{
			Type: "Bombola aggiuntiva", // Recupera 5 delle carte scartate
			Effects: []itemEffect{
				{
					effectType: RecoverDiscardedCards,
					value:      5,
				},
			},
		},
		{
			Type: "Sonar", // Guarda le prossime 5 carte ossigeno e mettile nell'ordine che vuoi
			Effects: []itemEffect{
				{
					effectType: LookAndReorder,
					value:      5,
				},
			},
		},
		{
			Type: "10  monete", // tesoro
			Effects: []itemEffect{
				{
					effectType: Treasure,
					value:      10,
				},
			},
		},
	}

	legendaryItems := []item{
		{
			Type: "Adrenalina", // riduce tutti i tipi di panico di 1
			Effects: []itemEffect{
				{
					effectType: ReduceAllPanicTypes,
					value:      1,
				},
			},
		},
	}

	/*
		specialItems := []string{
			"Amuleto DJ", // Serve a proteggersi da DJ
		}
	*/

	itemsDeck := make([]item, 0)
	for i := 0; i < commonItemsCount; i++ {
		itemsDeck = append(itemsDeck, commonItems[i])
	}

	for i := 0; i < uncommonItemsCount; i++ {
		itemsDeck = append(itemsDeck, uncommonItems[i])
	}

	for i := 0; i < rareItemsCount; i++ {
		itemsDeck = append(itemsDeck, rareItems[i])
	}

	for i := 0; i < veryRareItemsCount; i++ {
		itemsDeck = append(itemsDeck, veryRareItems[i])
	}

	for i := 0; i < legendaryItemsCount; i++ {
		itemsDeck = append(itemsDeck, legendaryItems[i])
	}

	o2Deck := make([]O2, 40)
	for i := range len(o2Deck) {

		o2 := O2{}

		switch (i / 10) + 1 {
		case 1:
			o2.Type = Encounter
			o2.Value = i + 1
			o2.ItemReward = 1
		case 2:
			o2.Type = Environment
			o2.Value = (i - 10) + 1
			o2.ItemReward = 1
		case 3:
			o2.Type = Technical
			o2.Value = (i - 20) + 1
			o2.ItemReward = 1
		case 4:
			o2.Type = Soprannatural
			o2.Value = (i - 30) + 1
			o2.ItemReward = 1
		}

		o2Deck[i] = o2

	}

	gamesWon := make(map[string]int)

	players := make([]player, NumberOfPlayers)
	for i := range NumberOfPlayers {
		player := player{
			Id: fmt.Sprintf("P%d", i+1),
			Ability: map[abilityType]abilityValue{
				Encounter:     {Value: 3, Modifier: 0},
				Environment:   {Value: 3, Modifier: 0},
				Technical:     {Value: 3, Modifier: 0},
				Soprannatural: {Value: 3, Modifier: 0},
			},
			Panic: map[abilityType]int{
				Encounter:     0,
				Environment:   0,
				Technical:     0,
				Soprannatural: 0,
			},
			PanicTollerance: map[abilityType]int{
				Encounter:     4,
				Environment:   4,
				Technical:     4,
				Soprannatural: 4,
			},
			O2:               make([]O2, len(o2Deck)),
			Inventory:        make([]item, 0),
			MaxInventorySize: 3,
		}

		gamesWon[player.Id] = 0

		copy(player.O2, o2Deck)

		randomizer.Shuffle(len(player.O2), func(i, j int) {
			player.O2[i], player.O2[j] = player.O2[j], player.O2[i]
		})

		players[i] = player
	}

	for game := range NumberOfGames {

		gameItems := make([]item, len(itemsDeck))

		copy(gameItems, itemsDeck)

		randomizer.Shuffle(len(gameItems), func(i, j int) {
			gameItems[i], gameItems[j] = gameItems[j], gameItems[i]
		})

		for p := range players {
			players[p].O2 = make([]O2, len(o2Deck))

			players[p].Panic = map[abilityType]int{
				Encounter:     0,
				Environment:   0,
				Technical:     0,
				Soprannatural: 0,
			}

			// Reset ability modifiers and effects for new game
			players[p].Ability = map[abilityType]abilityValue{
				Encounter:     {Value: 3, Modifier: 0},
				Environment:   {Value: 3, Modifier: 0},
				Technical:     {Value: 3, Modifier: 0},
				Soprannatural: {Value: 3, Modifier: 0},
			}
			players[p].Effects = make([]itemEffect, 0)
			players[p].Inventory = make([]item, 0)
			players[p].Treasure = 0
			players[p].MaxInventorySize = 3
			players[p].DiscardedO2 = make([]O2, 0)

			copy(players[p].O2, o2Deck)

			randomizer.Shuffle(len(players[p].O2), func(i, j int) {
				players[p].O2[i], players[p].O2[j] = players[p].O2[j], players[p].O2[i]
			})
		}

		fmt.Printf("%s\n", bold(cyan(fmt.Sprintf("START GAME %d", game+1))))

		round := 1
		for playersAlive(players) > 1 && len(gameItems) > 0 {

			fmt.Printf("\t%s\n", cyan(fmt.Sprintf("START ROUND %d", round)))

			for i := range players {

				if len(gameItems) == 0 {
					break
				}

				fmt.Printf("\t\t%s\n", bold(cyan(fmt.Sprintf("START TURN FOR %s", players[i].Id))))

				//BREATH
				cards, ok := drawO2CardForBreath(&players[i])
				if !ok {
					fmt.Printf("\t\t%s\n", red(bold(fmt.Sprintf("*** DEAD %s! ***", players[i].Id))))
					continue
				}

				diceResult := throwDice(4, randomizer)
				ability := players[i].Ability[cards[0].Type]
				if checkCardResolved(&players[i], cards[0], diceResult) {
					itemsToDraw := calculateItemsToDraw(cards[0].ItemReward, &players[i])
					items, newItems := draw(itemsToDraw, gameItems)
					gameItems = newItems

					// Process items: treasures go to treasure, others to inventory
					for _, item := range items {
						// Check if item is a treasure
						isTreasure := false
						treasureValue := 0
						for _, effect := range item.Effects {
							if effect.effectType == Treasure {
								isTreasure = true
								treasureValue += effect.value
							}
						}

						if isTreasure {
							// Treasure items go directly to treasure, not inventory
							players[i].Treasure += treasureValue
							fmt.Printf("\t\t\tBREATH '%s' %s %d + %d + %d > %d: found treasure = %s %s\n",
								cards[0].Type, green("RESOLVED"), ability.Value, ability.Modifier, diceResult, cards[0].Value,
								yellow(item.Type), green(fmt.Sprintf("(+%d treasure)", treasureValue)))
						} else {
							// Regular items go to inventory
							if len(players[i].Inventory) < players[i].MaxInventorySize {
								players[i].Inventory = append(players[i].Inventory, item)
								fmt.Printf("\t\t\tBREATH '%s' %s %d + %d + %d > %d: found item = %s\n",
									cards[0].Type, green("RESOLVED"), ability.Value, ability.Modifier, diceResult, cards[0].Value,
									cyan(item.Type))
							} else {
								fmt.Printf("\t\t\tBREATH '%s' %s %d + %d + %d > %d: found item = %s %s\n",
									cards[0].Type, green("RESOLVED"), ability.Value, ability.Modifier, diceResult, cards[0].Value,
									cyan(item.Type), yellow("(INVENTORY FULL, ITEM LOST)"))
							}
						}
						reducePanic(&players[i], cards[0].Type)
					}
				} else {
					if increasePanic(&players[i], cards[0].Type) {
						handlePanicTrigger(&players[i], cards[0].Type)
						continue
					} else {
						fmt.Printf("\t\t\tBREATH '%s' %s %d + %d + %d < %d: panic = %s\n",
							cards[0].Type, red("NOT RESOLVED"), ability.Value, ability.Modifier, diceResult, cards[0].Value,
							red(fmt.Sprintf("%d", players[i].Panic[cards[0].Type])))
					}
				}

				//EXECUTE ACTION
				for a := 1; a <= 3; a++ {
					if len(gameItems) == 0 {
						break
					}

					actionCompleted := false
					stopActions := false
					inputOk := false
					for !inputOk || !actionCompleted {
						// Print player status before action choice
						printPlayerStatus(&players[i])

						var input string
						fmt.Print("\t\t\tCHOOSE ACTION: (E)xplore - (U)se item - (P)ass")
						if a == 1 {
							fmt.Print(" - (C)alm down")
						}
						fmt.Print(": ")
						fmt.Scanln(&input)

						switch strings.ToLower(input) {
						case "c":
							if a != 1 {
								fmt.Printf("\t\t\t%s\n", yellow("CALM DOWN can only be used as the first action of the round"))
								continue
							}
							inputOk = true
							// Choose panic type
							panicTypeChosen := false
							var chosenPanicType abilityType
							for !panicTypeChosen {
								for panicType, panicValue := range players[i].Panic {
									if panicTypeChosen {
										break
									}
									var input string
									fmt.Printf("\t\t\tCHOOSE PANIC TYPE to calm down: %s (current: %d)? Y/N: ", panicType, panicValue)
									fmt.Scanln(&input)

									switch strings.ToLower(input) {
									case "y":
										chosenPanicType = panicType
										panicTypeChosen = true
									case "n":
										// Continue to next panic type
									default:
										fmt.Printf("\t\t\tINVALID INPUT\n")
										// Continue to next panic type
									}
									if panicTypeChosen {
										break
									}
								}
							}

							// Roll dice
							diceResult := throwDice(4, randomizer)
							currentPanic := players[i].Panic[chosenPanicType]

							fmt.Printf("\t\t\t%s: Rolled %s, Panic level: %s\n",
								cyan("CALM DOWN"), yellow(fmt.Sprintf("%d", diceResult)), cyan(fmt.Sprintf("%d", currentPanic)))

							if diceResult >= currentPanic {
								// Success: reduce panic by 1
								reducePanicBy(&players[i], chosenPanicType, 1)
								fmt.Printf("\t\t\t%s: %s %s panic reduced by 1 (now: %s)\n",
									cyan("CALM DOWN"), green("SUCCESS"), chosenPanicType, green(fmt.Sprintf("%d", players[i].Panic[chosenPanicType])))
							} else {
								// Failure: increase panic by 1 and check if panic triggers
								fmt.Printf("\t\t\t%s: %s %s panic will increase by 1\n",
									cyan("CALM DOWN"), red("FAILED"), chosenPanicType)
								if increasePanic(&players[i], chosenPanicType) {
									handlePanicTrigger(&players[i], chosenPanicType)
									stopActions = true
								} else {
									fmt.Printf("\t\t\t%s: %s %s panic increased by 1 (now: %s)\n",
										cyan("CALM DOWN"), red("FAILED"), chosenPanicType, red(fmt.Sprintf("%d", players[i].Panic[chosenPanicType])))
								}
							}

							actionCompleted = true
							// Calm down ends the player's turn
							stopActions = true
						case "u":
							if len(players[i].Inventory) == 0 {
								fmt.Printf("\t\t\tNO ITEMS IN INVENTORY\n")
								continue
							}
							choosen := false
							itemIndex := -1
							for !choosen {
								for idx, item := range players[i].Inventory {

									if choosen {
										break
									}

									fmt.Printf("\t\t\tCHOOSE ITEM: %s? Y/N: ", item.Type)
									fmt.Scanln(&input)

									switch strings.ToLower(input) {
									case "y":
										choosen = true
										itemIndex = idx
										for _, effect := range item.Effects {
											fmt.Printf("\t\t\tEFFECT: %s - Value: %d\n", effect.effectType, effect.value)
											switch effect.effectType {
											case IncreaseEncounter:
												increaseAbilityModifier(&players[i], Encounter, effect.value)
											case IncreaseEnvironment:
												increaseAbilityModifier(&players[i], Environment, effect.value)
											case IncreaseTechnical:
												increaseAbilityModifier(&players[i], Technical, effect.value)
											case IncreaseSoprannatural:
												increaseAbilityModifier(&players[i], Soprannatural, effect.value)
											case DrawMoreItems:
												players[i].Effects = append(players[i].Effects, effect)
											case Treasure:
												players[i].Treasure += effect.value
											case AddSlotInventory:
												players[i].MaxInventorySize += 1
											case RecoverDiscardedCards:
												cards, newDiscardedO2 := draw(effect.value, players[i].DiscardedO2)
												players[i].DiscardedO2 = newDiscardedO2
												players[i].O2 = append(cards, players[i].O2...)
											case FreeBreath:
												players[i].Effects = append(players[i].Effects, effect)
											case ReduceOnePanicType:
												choosen := false
												for !choosen {
													for ability, value := range players[i].Panic {

														if choosen {
															break
														}

														var input string
														fmt.Printf("\t\t\tREDUCE PANIC TYPE: %s = %d? Y/N: ", ability, value)
														fmt.Scanln(&input)

														switch strings.ToLower(input) {
														case "y":
															reducePanicBy(&players[i], ability, effect.value)
															choosen = true
														case "n":
															choosen = false
														default:
															fmt.Printf("\t\t\tINVALID INPUT")
															continue
														}
													}
												}
											case LookAndReorder:
												cards, newO2Cards := draw(effect.value, players[i].O2)

												sort.Slice(cards, func(i, j int) bool {
													return cards[i].Value < cards[j].Value
												})

												players[i].O2 = append(cards, newO2Cards...)

											case ReduceAllPanicTypes:
												for ability := range players[i].Panic {
													reducePanicBy(&players[i], ability, effect.value)
												}
											}
										}
										// Remove item from inventory after use
										if itemIndex >= 0 && itemIndex < len(players[i].Inventory) {
											players[i].Inventory = append(players[i].Inventory[:itemIndex], players[i].Inventory[itemIndex+1:]...)
										}
										inputOk = true
										// Item used, mark action as done
										actionCompleted = true
									case "n":
										continue
									default:
										fmt.Printf("\t\t\tINVALID INPUT\n")
										continue
									}
								}
							}
						case "e":
							inputOk = true
							cards, newO2 := draw(1, players[i].O2)
							if len(cards) == 0 {
								actionCompleted = true
								stopActions = true
								break
							}

							players[i].O2 = newO2
							players[i].DiscardedO2 = append(players[i].DiscardedO2, cards...)
							diceResult := throwDice(4, randomizer)
							ability := players[i].Ability[cards[0].Type]
							if checkCardResolved(&players[i], cards[0], diceResult) {
								itemsToDraw := calculateItemsToDraw(cards[0].ItemReward, &players[i])
								items, newItems := draw(itemsToDraw, gameItems)
								if len(items) == 0 {
									actionCompleted = true
									stopActions = true
								} else {
									addItemsToInventory(&players[i], items, "ACTION", cards[0].Type, ability.Value, ability.Modifier, diceResult, cards[0].Value)
									gameItems = newItems
									reducePanic(&players[i], cards[0].Type)
									// Action completed successfully, mark as done
									actionCompleted = true
								}
							} else {
								if increasePanic(&players[i], cards[0].Type) {
									handlePanicTrigger(&players[i], cards[0].Type)
									actionCompleted = true
									stopActions = true
								} else {
									fmt.Printf("\t\t\tACTION '%s' %s %d + %d + %d <= %d : panic = %s\n",
										cards[0].Type, red("NOT RESOLVED"), ability.Value, ability.Modifier, diceResult, cards[0].Value,
										red(fmt.Sprintf("%d", players[i].Panic[cards[0].Type])))
									// Action failed but was attempted, mark as done
									actionCompleted = true
								}
							}
						case "p":
							inputOk = true
							actionCompleted = true
							stopActions = true
						default:
							fmt.Printf("\t\t\tINVALID INPUT\n")
							continue
						}
					}

					if stopActions {
						break
					}

				}

				// Reset temporary effects at the end of the turn
				resetTemporaryEffects(&players[i])

				fmt.Printf("\t\t%s\n", cyan(fmt.Sprintf("END TURN FOR %s", players[i].Id)))

			}

			fmt.Printf("\t%s\n", cyan(fmt.Sprintf("END ROUND %d", round)))
			round++
		}

		fmt.Printf("%s\n", cyan(fmt.Sprintf("ITEMS LEFT: %d", len(gameItems))))

		fmt.Printf("%s\n", bold(cyan(fmt.Sprintf("END GAME %d", game+1))))

	}

	fmt.Println(gamesWon)
}

func playersAlive(players []player) int {

	alivePlayers := len(players)
	for _, player := range players {
		if len(player.O2) == 0 {
			alivePlayers--
		}

	}

	return alivePlayers

}

func draw[T any](numberOfElementsToDraw int, slice []T) ([]T, []T) {

	if len(slice) <= numberOfElementsToDraw {
		numberOfElementsToDraw = len(slice)
	}

	drawedElements := make([]T, numberOfElementsToDraw)
	for i := 0; i < numberOfElementsToDraw; i++ {
		drawedElement := slice[0]
		drawedElements[i] = drawedElement
		slice = slice[1:]
	}

	return drawedElements, slice
}

func look[T any](numberOfElementsToLook int, slice []T) []T {

	if len(slice) <= numberOfElementsToLook {
		numberOfElementsToLook = len(slice)
	}

	drawedElements := make([]T, numberOfElementsToLook)
	for i := 0; i < numberOfElementsToLook; i++ {
		drawedElement := slice[0]
		drawedElements[i] = drawedElement
	}

	return drawedElements
}

func throwDice(faces int, randomizer *rand.Rand) int {
	return randomizer.Intn(faces) + 1
}

/*

I giocatori possono muoversi su 3 livelli di profondità
Al livello 4 c'è DJ
I giocatori per fare azioni consumano ossigeno pescando carte da un mazzo che ne rappresenta la riserva.
Finite le carte ossigeno il giocatore muore
Le carte pescate possono far aumentare il livello di panico del giocatore.
Quando il livello di panico supera una certa soglia si attiva un effetto molto negativo
Esplorando i fondali è possibile trovare oggetti e tesori

*/

/*

Giocatore parte con:
 - Equip base:
	- Pinne
	- Bombola
	- Maschera

	Avanzato: si possono spendere soldi per partire con equip extra

Resistenze base:
	- Fisico
	- Mentale
	- Empatia
	- Superstizione

	Avanzato: si possono modificare i valori aggiungendo o sottraendo punti mantenendo però il totale invariato

1 Mazzo oggetti comune a tutti i giocatori:

- Carte oggetti utilità
- Carte tesoro
- Carte amuleto

1 Mazzo ossigno per giocatore:

Ogni carta ha delle icone e un numero che ne rappresenta il valore
- Carte fanno accumulare un tipo di panico (superata una soglia scatta il panico)
- Il panico si accumula quando provando a fare un'azione fallisci
- I giocatori saranno più bravi a fare certe cose e meno bravi a farne altre (vedi resistenze)

Ogni turno il giocatore dovrà pescare delle carte respiro.
Queste carte avranno un numero e delle icone
	- Tridente: rappresenta la capacità di fare azioni offensive
	- Sonar: capacità esplorativa
	- Esoterica: capacità di r

	In base alle icone e ai punti azione trovati il giocatore dovrà decidere cosa fare
	Se riuscirà nell'azione la risolverà
	Altrimenti aumenterà il panico nella caratteristica corrispondente
	Es. se attacca qualcuno o qualcosa e fallisce allora aumenterà il panico fisico
		se esplorerà e fallirà nell'esplorazione allora aumenterà il panico mentale

	Superata la soglia di panico per quella caratteristica il giocatore subità gli effetti in base al tipo di panico:
		- Empatia scatenerà paranoia e vorrà attaccare gli altri etc.
*/

/*

Loop di gioco

	1. Respiro: Pesco una carta ossigeno e provo a risolverla
				Se la risolvo con successo guadagno o punti per le prossime azioni oppure oggetti
				Se non la risolvo prendo punti panico

				Se supero una certa dose di panico si attiva l'effetto

	2. Scelta dell'azione costo 1 O2:
		- Mi sposto di un livello in basso
		- Risalgo di un livello
		- Esploro

		Provo a risolvere l'azione (vedi respiro)

		Se ho fatto meno di 3 azioni posso continuare e provare a farne altre.



*/
