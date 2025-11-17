package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var NumberOfGames = 1
var NumberOfPlayers = 2

type O2 struct {
	Type       string
	Value      int
	ItemReward int
}

type player struct {
	Id              string
	O2              []O2
	Ability         map[string]int
	Panic           map[string]int
	PanicTollerance map[string]int
	RiskTollerance  int
	Treasure        int
	Inventory       []item
}

type item struct {
	Type string
}

type objective struct {
	Type         string
	GameDuration int
	WinCondition []itemCondition
}

type itemCondition struct {
	Item  string
	Count int
}

func main() {

	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	commonItemsCount := 8
	uncommonItemsCount := 8
	rareItemsCount := 4
	veryRareItemsCount := 2
	legendaryItemsCount := 1

	commonItems := []string{
		"Coltello",            // +2 agli incontri
		"Torcia",              // +2 alle esplorazioni
		"Barometro",           // +2 agli imprevisti ambientali
		"1 moneta",            // tesoro
		"Sacca maggiorata",    // aggiunge uno slot in più all'inventario
		"Maschera migliorata", // +1 alle esplorazioni
		"Pinne migliorate",    // +1 agli imprevisti ambientali
		"Bombola migliorata",  // recupera una carta ossigeno scartata
	}

	uncommonItems := []string{
		"Coltello migliorato",    // +3 agli incontri
		"Torcia migliorata",      // +3 alle esplorazioni
		"Barometro migliorato",   // +3 agli imprevisti ambientali
		"Fiocina",                // +4 agli incontri
		"2 monete",               // tesoro
		"Bombola maggiorata",     // recupera 2 carte ossigeno scartate
		"amuleto di poseidone",   // +2 agli incontri soprannaturali
		"Respiratore migliorato", // respira gratis
	}

	rareItems := []string{
		"Bombola aggiuntiva",    // recupera 3 carte ossigeno scartate
		"tridente di poseidone", // +2 agli incontri e agli imprevisti sovrannaturali
		"Kit antistress",        // riduce di 2 un tipo di panico
		"3 monete",              // tesoro
	}

	veryRareItems := []string{
		"Sonar",    // Guarda le prossime 5 carte ossigeno e mettile nell'ordine che vuoi
		"5 monete", // tesoro
	}

	legendaryItems := []string{
		"Adrenalina", // Rirudce di 1 tutti i tipi di panico
	}

	/*
		specialItems := []string{
			"Amuleto DJ", // Serve a proteggersi da DJ
		}
	*/

	itemsDeck := make([]item, 0)
	for i := 0; i < commonItemsCount; i++ {

		itemsDeck = append(itemsDeck, item{
			Type: commonItems[i],
		})

		itemsDeck = append(itemsDeck, item{
			Type: commonItems[i],
		})

	}

	for i := 0; i < uncommonItemsCount; i++ {
		itemsDeck = append(itemsDeck, item{
			Type: uncommonItems[i],
		})
	}

	for i := 0; i < rareItemsCount; i++ {
		itemsDeck = append(itemsDeck, item{
			Type: rareItems[i],
		})
	}

	for i := 0; i < veryRareItemsCount; i++ {
		itemsDeck = append(itemsDeck, item{
			Type: veryRareItems[i],
		})
	}

	for i := 0; i < legendaryItemsCount; i++ {
		itemsDeck = append(itemsDeck, item{
			Type: legendaryItems[i],
		})
	}

	o2Deck := make([]O2, 40)
	for i := range len(o2Deck) {

		o2 := O2{}

		switch (i / 10) + 1 {
		case 1:
			o2.Type = "ENCOUNTER"
			o2.Value = i + 1
			o2.ItemReward = 1
		case 2:
			o2.Type = "ENVIRONMENT"
			o2.Value = (i - 10) + 1
			o2.ItemReward = 1
		case 3:
			o2.Type = "TECHNICAL"
			o2.Value = (i - 20) + 1
			o2.ItemReward = 1
		case 4:
			o2.Type = "SOPRANNATURAL"
			o2.Value = (i - 30) + 1
			o2.ItemReward = 1
		}

		o2Deck[i] = o2

	}

	gamesWon := make(map[string]int)

	players := make([]player, NumberOfPlayers)
	for i := range NumberOfPlayers {
		player := player{
			Id:             fmt.Sprintf("P%d", i+1),
			Treasure:       0,
			RiskTollerance: (i / 2) + 1,
			Ability: map[string]int{
				"ENCOUNTER":     3,
				"ENVIRONMENT":   3,
				"TECHNICAL":     3,
				"SOPRANNATURAL": 3,
			},
			Panic: map[string]int{
				"ENCOUNTER":     0,
				"ENVIRONMENT":   0,
				"TECHNICAL":     0,
				"SOPRANNATURAL": 0,
			},
			PanicTollerance: map[string]int{
				"ENCOUNTER":     3,
				"ENVIRONMENT":   3,
				"TECHNICAL":     3,
				"SOPRANNATURAL": 3,
			},
			O2: make([]O2, len(o2Deck)),
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
			players[p].Treasure = 0

			players[p].Panic = map[string]int{
				"ENCOUNTER":     0,
				"ENVIRONMENT":   0,
				"TECHNICAL":     0,
				"SOPRANNATURAL": 0,
			}

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
				cards, newO2 := draw(1, players[i].O2)
				if len(cards) == 0 {
					fmt.Printf("\t\t*** DEAD %s! ***\n", players[i].Id)
					continue
				}
				players[i].O2 = newO2
				diceResult := throwDice(4, randomizer)
				if players[i].Ability[cards[0].Type]+diceResult > cards[0].Value {
					items, newItems := draw(cards[0].ItemReward, gameItems)
					gameItems = newItems

					for _, item := range items {
						fmt.Printf("\t\t\tBREATH RESOLVED %d + %d > %d: found item = %s\n", players[i].Ability[cards[0].Type], diceResult, cards[0].Value, item.Type)
						if players[i].Panic[cards[0].Type] > 0 {
							players[i].Panic[cards[0].Type] -= 1
						}
					}
				} else {
					players[i].Panic[cards[0].Type] += 1
					fmt.Printf("\t\t\tBREATH NOT RESOLVED %d + %d < %d: panic = %d\n", players[i].Ability[cards[0].Type], diceResult, cards[0].Value, players[i].Panic[cards[0].Type])
					if players[i].Panic[cards[0].Type] == players[i].PanicTollerance[cards[0].Type] {
						fmt.Printf("\t\t\t*** PANIC! ***\n")
						players[i].Panic[cards[0].Type] = 0
						players[i].Treasure = 0
						continue
					}
				}

				//TAKE ACTION
				for a := 1; a <= 3; a++ {
					if len(gameItems) == 0 {
						break
					}

					pass := false
					inputOk := false
					for !inputOk || !pass {
						var input string
						fmt.Printf("\t\t\tCHOOSE ACTION: (E)xplore - (P)ass: \n")
						fmt.Scanln(&input)

						switch strings.ToLower(input) {
						case "e":
							inputOk = true
							cards, newO2 := draw(1, players[i].O2)
							if len(cards) == 0 {
								pass = true

							}

							players[i].O2 = newO2
							diceResult := throwDice(4, randomizer)
							if players[i].Ability[cards[0].Type]+diceResult > cards[0].Value {
								items, newItems := draw(cards[0].ItemReward, gameItems)
								if len(items) == 0 {
									pass = true
								}
								gameItems = newItems

								for _, item := range items {
									fmt.Printf("\t\t\tACTION '%s' RESOLVED %d + %d > %d: item found = %s\n", cards[0].Type, players[i].Ability[cards[0].Type], diceResult, cards[0].Value, item.Type)
								}
								if players[i].Panic[cards[0].Type] > 0 {
									players[i].Panic[cards[0].Type] -= 1
								}
							} else {
								players[i].Panic[cards[0].Type] += 1
								fmt.Printf("\t\t\tACTION '%s' NOT RESOLVED %d + %d <= %d : panic = %d\n", cards[0].Type, players[i].Ability[cards[0].Type], diceResult, cards[0].Value, players[i].Panic[cards[0].Type])
								if players[i].Panic[cards[0].Type] == players[i].PanicTollerance[cards[0].Type] {
									fmt.Printf("\t\t\t*** PANIC! ***\n")
									players[i].Panic[cards[0].Type] = 0
									players[i].Treasure = 0
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

func winner(players []player, objective objective) string {
	switch objective.Type {
	case "MAX_TREASURES":
		bestTreasure := 0
		playerId := ""
		for _, player := range players {
			if player.Treasure > bestTreasure {
				bestTreasure = player.Treasure
				playerId = player.Id
			}
		}

		return playerId
	case "SURVIVE_DJ":

	case "KILL_DY":
	}

	return ""
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
