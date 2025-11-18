package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

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
				Encounter:     3,
				Environment:   3,
				Technical:     3,
				Soprannatural: 3,
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

		fmt.Printf("START GAME %d\n", game+1)

		round := 1
		for playersAlive(players) > 1 && len(gameItems) > 0 {

			fmt.Printf("\tSTART ROUND %d\n", round)

			for i := range players {

				if len(gameItems) == 0 {
					break
				}

				fmt.Printf("\t\tSTART TURN FOR %s\n", players[i].Id)

				//BREATH
				// Check for FreeBreath effect
				hasFreeBreath := false
				for _, effect := range players[i].Effects {
					if effect.effectType == FreeBreath {
						hasFreeBreath = true
						break
					}
				}

				var cards []O2
				var newO2 []O2
				if !hasFreeBreath {
					cards, newO2 = draw(1, players[i].O2)
					if len(cards) == 0 {
						fmt.Printf("\t\t*** DEAD %s! ***\n", players[i].Id)
						continue
					}
					players[i].O2 = newO2
					players[i].DiscardedO2 = append(players[i].DiscardedO2, cards...)
				} else {
					// FreeBreath: draw card but put it back in the deck instead of discarding
					cards, newO2 = draw(1, players[i].O2)
					if len(cards) == 0 {
						fmt.Printf("\t\t*** DEAD %s! ***\n", players[i].Id)
						continue
					}
					players[i].O2 = newO2
					// Card is put back at the bottom of the deck with FreeBreath
					players[i].O2 = append(players[i].O2, cards[0])
				}

				diceResult := throwDice(4, randomizer)
				if players[i].Ability[cards[0].Type].Value+players[i].Ability[cards[0].Type].Modifier+diceResult > cards[0].Value {
					// Calculate number of items to draw (including DrawMoreItems effects)
					itemsToDraw := cards[0].ItemReward
					for _, effect := range players[i].Effects {
						if effect.effectType == DrawMoreItems {
							itemsToDraw += effect.value
						}
					}

					items, newItems := draw(itemsToDraw, gameItems)
					gameItems = newItems

					for _, item := range items {
						fmt.Printf("\t\t\tBREATH '%s' RESOLVED %d + %d + %d > %d: found item = %s\n", cards[0].Type, players[i].Ability[cards[0].Type].Value, players[i].Ability[cards[0].Type].Modifier, diceResult, cards[0].Value, item.Type)
						if players[i].Panic[cards[0].Type] > 0 {
							players[i].Panic[cards[0].Type] -= 1
							if players[i].Panic[cards[0].Type] < 0 {
								players[i].Panic[cards[0].Type] = 0
							}
						}
					}
				} else {
					players[i].Panic[cards[0].Type] += 1
					fmt.Printf("\t\t\tBREATH '%s' NOT RESOLVED %d + %d + %d < %d: panic = %d\n", cards[0].Type, players[i].Ability[cards[0].Type].Value, players[i].Ability[cards[0].Type].Modifier, diceResult, cards[0].Value, players[i].Panic[cards[0].Type])
					if players[i].Panic[cards[0].Type] >= players[i].PanicTollerance[cards[0].Type] {
						fmt.Printf("\t\t\t*** PANIC! ***\n")
						players[i].Panic[cards[0].Type] = 0
						continue
					}
				}

				//EXECUTE ACTION
				for a := 1; a <= 3; a++ {
					if len(gameItems) == 0 {
						break
					}

					pass := false
					inputOk := false
					for !inputOk || !pass {
						var input string
						fmt.Printf("\t\t\tCHOOSE ACTION: (E)xplore - (U)se item - (P)ass: ")
						fmt.Scanln(&input)

						switch strings.ToLower(input) {
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

									fmt.Printf("\t\t\tCHOOSE ITEM: %s? Y/N", item.Type)
									fmt.Scanln(&input)

									switch strings.ToLower(input) {
									case "y":
										choosen = true
										itemIndex = idx
										for _, effect := range item.Effects {
											fmt.Printf("\t\t\tEFFECT: %s - Value: %d\n", effect.effectType, effect.value)
											switch effect.effectType {
											case IncreaseEncounter:
												ability := players[i].Ability[Encounter]
												ability.Modifier += effect.value
												players[i].Ability[Encounter] = ability
											case IncreaseEnvironment:
												ability := players[i].Ability[Environment]
												ability.Modifier += effect.value
												players[i].Ability[Environment] = ability
											case IncreaseTechnical:
												ability := players[i].Ability[Technical]
												ability.Modifier += effect.value
												players[i].Ability[Technical] = ability
											case IncreaseSoprannatural:
												ability := players[i].Ability[Soprannatural]
												ability.Modifier += effect.value
												players[i].Ability[Soprannatural] = ability
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
														fmt.Printf("\t\t\tREDuCE PANIC TYPE: %s = %d? Y/N", ability, value)
														fmt.Scanln(&input)

														switch strings.ToLower(input) {
														case "y":
															players[i].Panic[ability] = players[i].Panic[ability] - effect.value
															if players[i].Panic[ability] < 0 {
																players[i].Panic[ability] = 0
															}
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
													players[i].Panic[ability] = players[i].Panic[ability] - effect.value
													if players[i].Panic[ability] < 0 {
														players[i].Panic[ability] = 0
													}
												}
											}
										}
										// Remove item from inventory after use
										if itemIndex >= 0 && itemIndex < len(players[i].Inventory) {
											players[i].Inventory = append(players[i].Inventory[:itemIndex], players[i].Inventory[itemIndex+1:]...)
										}
										inputOk = true
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
								pass = true
								break
							}

							players[i].O2 = newO2
							players[i].DiscardedO2 = append(players[i].DiscardedO2, cards...)
							diceResult := throwDice(4, randomizer)
							if players[i].Ability[cards[0].Type].Value+players[i].Ability[cards[0].Type].Modifier+diceResult > cards[0].Value {
								// Calculate number of items to draw (including DrawMoreItems effects)
								itemsToDraw := cards[0].ItemReward
								for _, effect := range players[i].Effects {
									if effect.effectType == DrawMoreItems {
										itemsToDraw += effect.value
									}
								}

								items, newItems := draw(itemsToDraw, gameItems)
								if len(items) == 0 {
									pass = true
								}

								// Add items to inventory if there's space
								itemsAdded := 0
								for _, item := range items {
									if len(players[i].Inventory) < players[i].MaxInventorySize {
										players[i].Inventory = append(players[i].Inventory, item)
										itemsAdded++
										fmt.Printf("\t\t\tACTION '%s' RESOLVED %d + %d + %d > %d: item found = %s\n", cards[0].Type, players[i].Ability[cards[0].Type].Value, players[i].Ability[cards[0].Type].Modifier, diceResult, cards[0].Value, item.Type)
									} else {
										fmt.Printf("\t\t\tACTION '%s' RESOLVED %d + %d + %d > %d: item found = %s (INVENTORY FULL, ITEM LOST)\n", cards[0].Type, players[i].Ability[cards[0].Type].Value, players[i].Ability[cards[0].Type].Modifier, diceResult, cards[0].Value, item.Type)
									}
								}
								gameItems = newItems

								if players[i].Panic[cards[0].Type] > 0 {
									players[i].Panic[cards[0].Type] -= 1
									if players[i].Panic[cards[0].Type] < 0 {
										players[i].Panic[cards[0].Type] = 0
									}
								}
							} else {
								players[i].Panic[cards[0].Type] += 1
								fmt.Printf("\t\t\tACTION '%s' NOT RESOLVED %d + %d + %d <= %d : panic = %d\n", cards[0].Type, players[i].Ability[cards[0].Type].Value, players[i].Ability[cards[0].Type].Modifier, diceResult, cards[0].Value, players[i].Panic[cards[0].Type])
								if players[i].Panic[cards[0].Type] >= players[i].PanicTollerance[cards[0].Type] {
									fmt.Printf("\t\t\t*** PANIC! ***\n")
									players[i].Panic[cards[0].Type] = 0
									pass = true
								}
							}
						case "p":
							inputOk = true
							pass = true
						default:
							fmt.Printf("\t\t\tINVALID INPUT\n")
							continue
						}
					}

					if pass {
						break
					}

				}

				fmt.Printf("\t\tEND TURN FOR %s\n", players[i].Id)

			}

			fmt.Printf("\tEND ROUND %d\n", round)
			round++
		}

		fmt.Printf("ITEMS LEFT: %d\n", len(gameItems))

		fmt.Printf("END GAME %d\n", game+1)

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
