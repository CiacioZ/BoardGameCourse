package main

import (
	"fmt"
	"math/rand"
	"time"
)

type O2 struct {
	Type     string
	Value    int
	Treasure int
}

type player struct {
	Id              string
	O2              []O2
	Ability         int
	Treasure        int
	Panic           int
	PanicTollerance int
}

type object struct {
	Type  string
	Value int
}

var NumberOfGames = 1
var NumberOfPlayers = 2

func main() {

	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	itemsDeck := make([]object, 30+5*NumberOfPlayers)
	for i := range len(itemsDeck) {

		item := object{
			Type:  "TREASURE",
			Value: (i / 10) + 1,
		}

		itemsDeck[i] = item
	}

	o2Deck := make([]O2, 40)
	for i := range len(o2Deck) {

		o2 := O2{}

		switch (i / 10) + 1 {
		case 1:
			o2.Type = "ENCOUNTER"
			o2.Value = i + 1
			o2.Treasure = 1
		case 2:
			o2.Type = "EXPLORATION"
			o2.Value = (i - 10) + 1
			o2.Treasure = 1
		case 3:
			o2.Type = "UNEXPECTED"
			o2.Value = (i - 20) + 1
			o2.Treasure = 1
		case 4:
			o2.Type = "SOPRANNATURAL"
			o2.Value = (i - 30) + 1
			o2.Treasure = 1
		}

		o2Deck[i] = o2

	}

	players := make([]player, NumberOfPlayers)
	for i := range NumberOfPlayers {
		player := player{
			Id:              fmt.Sprintf("P%d", i+1),
			Treasure:        0,
			Ability:         3,
			PanicTollerance: 3,
			O2:              make([]O2, len(o2Deck)),
		}

		copy(player.O2, o2Deck)

		randomizer.Shuffle(len(player.O2), func(i, j int) {
			player.O2[i], player.O2[j] = player.O2[j], player.O2[i]
		})

		players[i] = player
	}

	for game := range NumberOfGames {

		gameItems := make([]object, len(itemsDeck))

		copy(gameItems, itemsDeck)

		randomizer.Shuffle(len(gameItems), func(i, j int) {
			gameItems[i], gameItems[j] = gameItems[j], gameItems[i]
		})

		fmt.Printf("START GAME %d\n", game)

		round := 1
		for playersAlive(players) > 1 {

			fmt.Printf("\tSTART ROUND %d\n", round)

			for i := range players {

				fmt.Printf("\t\tSTART TURN FOR %s\n", players[i].Id)

				//BREATH
				cards, newO2 := draw(1, players[i].O2)
				players[i].O2 = newO2
				diceResult := throwDice(4, randomizer)
				if players[i].Ability+diceResult >= cards[0].Value {
					items, newItems := draw(cards[0].Treasure, gameItems)
					gameItems = newItems

					for _, item := range items {
						players[i].Treasure += item.Value
						fmt.Printf("\t\t\tBREATH RESOLVED %d + %d >= %d: treasure = %d\n", players[i].Ability, diceResult, cards[0].Value, players[i].Treasure)
						if players[i].Panic > 0 {
							players[i].Panic -= 1
						}
					}
				} else {
					players[i].Panic += 1
					fmt.Printf("\t\t\tBREATH NOT RESOLVED %d + %d < %d: panic = %d\n", players[i].Ability, diceResult, cards[0].Value, players[i].Panic)
					if players[i].Panic == players[i].PanicTollerance {
						fmt.Printf("\t\t\tPANIC!\n")
						players[i].Panic = 0
						players[i].Treasure = 0
						continue
					}
				}

				//TAKE ACTION
				for a := 1; a <= 3; a++ {
					if players[i].Panic == players[i].PanicTollerance-1 {
						if a == 1 {
							diceResult := throwDice(4, randomizer)
							if diceResult >= players[i].Panic {
								fmt.Printf("\t\t\tACTION 'CALM_DOWN' SUCCESSFUL! %d >= %d\n", diceResult, players[i].Panic)
								players[i].Panic -= 1
							} else {
								fmt.Printf("\t\t\tACTION 'CALM_DOWN' FAILED! %d < %d\n", diceResult, players[i].Panic)
							}
							break
						} else {
							break
						}
					}
					cards, newO2 := draw(1, players[i].O2)
					if len(cards) == 0 {
						break
					}

					players[i].O2 = newO2
					if players[i].Ability+throwDice(4, randomizer) >= cards[0].Value {
						items, newItems := draw(cards[0].Treasure, gameItems)
						gameItems = newItems

						for _, item := range items {
							players[i].Treasure += item.Value
							fmt.Printf("\t\t\tACTION '%s' RESOLVED: treasure = %d\n", cards[0].Type, players[i].Treasure)
						}
						if players[i].Panic > 0 {
							players[i].Panic -= 1
						}
					} else {
						players[i].Panic += 1
						fmt.Printf("\t\t\tACTION '%s' NOT RESOLVED: panic = %d\n", cards[0].Type, players[i].Panic)
						if players[i].Panic == players[i].PanicTollerance {
							fmt.Printf("\t\t\tPANIC!\n")
							players[i].Panic = 0
							players[i].Treasure = 0
							break
						}
					}
				}

				fmt.Printf("\t\tEND TURN FOR %s\n", players[i].Id)

			}

			fmt.Printf("\tEND ROUND %d\n", round)
			round++
		}

		fmt.Printf("END GAME %d\n", game)

		winner, treasure := winner(players)

		fmt.Printf("WINNER: %s has %d treasure points\n", winner, treasure)
	}
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

func winner(players []player) (string, int) {
	bestTreasure := 0
	playerId := ""
	for _, player := range players {
		if player.Treasure > bestTreasure {
			bestTreasure = player.Treasure
			playerId = player.Id
		}
	}

	return playerId, bestTreasure
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
