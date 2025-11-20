package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// --- ANSI COLOR CODES ---
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

// --- COLOR HELPER FUNCTIONS ---
func colorize(text string, color string) string { return color + text + colorReset }
func red(text string) string                    { return colorize(text, colorRed) }
func green(text string) string                  { return colorize(text, colorGreen) }
func yellow(text string) string                 { return colorize(text, colorYellow) }
func cyan(text string) string                   { return colorize(text, colorCyan) }
func bold(text string) string                   { return colorize(text, colorBold) }
func purple(text string) string                 { return colorize(text, colorPurple) }

var NumberOfGames = 1
var NumberOfPlayers = 2

// --- EXCEL HELPER STRUCT ---
type GameLogger struct {
	File      *excelize.File
	SheetName string
	RowIndex  int
}

func NewGameLogger() *GameLogger {
	f := excelize.NewFile()
	return &GameLogger{
		File:     f,
		RowIndex: 1,
	}
}

func (gl *GameLogger) InitSheet(gameNum int) {
	sheetName := fmt.Sprintf("Game %d", gameNum)
	gl.SheetName = sheetName
	index, _ := gl.File.NewSheet(sheetName)
	gl.File.SetActiveSheet(index)
	gl.RowIndex = 1

	// Headers
	headers := []string{"Round", "Player", "Phase", "Action", "Details", "Outcome", "Panic", "HP (O2)", "Points", "Treasure", "Log Message"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		gl.File.SetCellValue(sheetName, cell, h)
	}
	gl.RowIndex++
}

// LogEvent scrive su Excel e (opzionalmente) su console se serve debugging extra,
// ma qui lo usiamo per tracciare i dati silenziosamente mentre il gioco stampa a video.
func (gl *GameLogger) LogEvent(round int, p *player, phase, action, details, outcome, message string) {
	if gl == nil {
		return
	}

	// Valori di default se p è nil (Evento di Sistema)
	playerId := "System"
	panicVal := 0
	o2Len := 0
	totalPoints := 0
	treasure := 0

	// Se il giocatore esiste, sovrascrivi con i dati reali
	if p != nil {
		playerId = p.Id
		panicVal = p.Panic
		o2Len = len(p.O2)
		treasure = p.Treasure
		for _, v := range p.AbilityPool {
			totalPoints += v
		}
	}

	// Scrive riga usando le variabili calcolate
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("A%d", gl.RowIndex), round)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("B%d", gl.RowIndex), playerId)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("C%d", gl.RowIndex), phase)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("D%d", gl.RowIndex), action)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("E%d", gl.RowIndex), details)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("F%d", gl.RowIndex), outcome)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("G%d", gl.RowIndex), panicVal)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("H%d", gl.RowIndex), o2Len)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("I%d", gl.RowIndex), totalPoints)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("J%d", gl.RowIndex), treasure)
	gl.File.SetCellValue(gl.SheetName, fmt.Sprintf("K%d", gl.RowIndex), message)

	gl.RowIndex++
}

func (gl *GameLogger) Save(filename string) {
	// Rimuove Sheet1 default se non usato
	gl.File.DeleteSheet("Sheet1")
	if err := gl.File.SaveAs(filename); err != nil {
		fmt.Println(red("Error saving Excel file:"), err)
	} else {
		fmt.Println(green(fmt.Sprintf("Game logs saved to %s", filename)))
	}
}

// Rimuove i codici colore ANSI per il log Excel
func stripAnsi(str string) string {
	// Implementazione molto basilare, rimuove il carattere Escape principale
	// Per una pulizia perfetta servirebbe una regex, ma per i log excel basta non avere troppa sporcizia
	// Qui salviamo il messaggio raw o una versione semplificata
	return str
}

// --- TYPES & CONSTANTS ---
type abilityType string

const (
	Encounter     abilityType = "ENCOUNTER"
	Environment   abilityType = "ENVIRONMENT"
	Technical     abilityType = "TECHNICAL"
	Soprannatural abilityType = "SOPRANNATURAL"
)

var allAbilities = []abilityType{Encounter, Environment, Technical, Soprannatural}

type O2 struct {
	Name          string
	Type          abilityType
	SecondaryType abilityType
	Value         int
}

type player struct {
	Id               string
	O2               []O2
	DiscardedO2      []O2
	AbilityPool      map[abilityType]int
	Panic            int
	Treasure         int
	Inventory        []item
	Effects          []itemEffect
	MaxInventorySize int
	RoundScore       int
}

type item struct {
	Type    string
	Effects []itemEffect
}

type effectType string

const (
	AddAbilityPoints      effectType = "ADD_ABILITY_POINTS"
	DrawMoreItems         effectType = "DRAW_MORE_ITEMS"
	Treasure              effectType = "TREASURE"
	AddSlotInventory      effectType = "ADD_SLOT_INVENTORY"
	RecoverDiscardedCards effectType = "RECOVER_DISCARDED_CARDS"
	FreeBreath            effectType = "FREE_BREATH"
	ReducePanic           effectType = "REDUCE_PANIC"
	LookAndReorder        effectType = "LOOK_AND_REORDER"
)

type itemEffect struct {
	effectType effectType
	targetType abilityType
	value      int
}

// --- HELPER FUNCTIONS ---

func getDiceFaces(panicLevel int) int {
	switch panicLevel {
	case 0:
		return 8
	case 1:
		return 6
	case 2:
		return 4
	default:
		return 0
	}
}

func throwExplodingDice(faces int, randomizer *rand.Rand) int {
	total := 0
	rollsStr := ""
	for {
		roll := randomizer.Intn(faces) + 1
		total += roll
		if rollsStr != "" {
			rollsStr += " + "
		}
		if roll == faces {
			rollsStr += bold(purple(fmt.Sprintf("%d(CRIT!)", roll)))
		} else {
			rollsStr += fmt.Sprintf("%d", roll)
			break
		}
	}
	if total > faces {
		fmt.Printf("\t\t\t(Rolled: %s = %d)\n", rollsStr, total)
	}
	return total
}

func hasFreeBreath(p *player) bool {
	for _, effect := range p.Effects {
		if effect.effectType == FreeBreath {
			return true
		}
	}
	return false
}

func handlePanicTrigger(p *player, logger *GameLogger, round int) {
	fmt.Printf("\t\t\t%s\n", red(bold("*** PANIC ATTACK! (Level 3 Reached) ***")))

	logger.LogEvent(round, p, "Panic", "Trigger", "Panic Level 3", "Disaster", "Lose Inventory & 5 O2")

	p.Panic = 0
	if len(p.Inventory) > 0 {
		fmt.Printf("\t\t\t%s\n", yellow("Lost all items from inventory due to panic!"))
		p.Inventory = make([]item, 0)
	}
	cardsToDiscard, newO2AfterPanic := draw(5, p.O2)
	p.O2 = newO2AfterPanic
	p.DiscardedO2 = append(p.DiscardedO2, cardsToDiscard...)
	fmt.Printf("\t\t\t%s\n", yellow(fmt.Sprintf("Gasping for air! Discarded %d O2 cards.", len(cardsToDiscard))))
}

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
		cards, newO2 = draw(1, p.O2)
		if len(cards) == 0 {
			return nil, false
		}
		p.O2 = newO2
		p.O2 = append(p.O2, cards[0])
	}
	return cards, true
}

func resolveCardInteraction(p *player, card O2, randomizer *rand.Rand, actionName string, logger *GameLogger, round int) bool {
	if p.Panic >= 3 {
		handlePanicTrigger(p, logger, round)
		return false
	}
	diceFaces := getDiceFaces(p.Panic)

	cardDesc := string(card.Type)
	if card.SecondaryType != "" {
		cardDesc += " + " + string(card.SecondaryType) + " (BOSS)"
	}

	availablePrimary := p.AbilityPool[card.Type]
	availableSecondary := 0
	if card.SecondaryType != "" {
		availableSecondary = p.AbilityPool[card.SecondaryType]
	}

	fmt.Printf("\t\t\t%s: %s (Difficulty: %d) | Panic: %d (Dice: d%d)\n",
		actionName, cyan(cardDesc), card.Value, p.Panic, diceFaces)
	fmt.Printf("\t\t\tPools -> %s: %d", card.Type, availablePrimary)
	if card.SecondaryType != "" {
		fmt.Printf(" | %s: %d", card.SecondaryType, availableSecondary)
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	spentPrimary := 0
	spentSecondary := 0

	for {
		fmt.Printf("\t\t\tSpend %s points? (0-%d): ", card.Type, availablePrimary)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		val, err := strconv.Atoi(input)
		if err == nil && val >= 0 && val <= availablePrimary {
			spentPrimary = val
			break
		}
		fmt.Println("\t\t\tInvalid amount.")
	}

	if card.SecondaryType != "" {
		for {
			fmt.Printf("\t\t\tSpend %s points? (0-%d): ", card.SecondaryType, availableSecondary)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			val, err := strconv.Atoi(input)
			if err == nil && val >= 0 && val <= availableSecondary {
				spentSecondary = val
				break
			}
			fmt.Println("\t\t\tInvalid amount.")
		}
	}

	p.AbilityPool[card.Type] -= spentPrimary
	if card.SecondaryType != "" {
		p.AbilityPool[card.SecondaryType] -= spentSecondary
	}

	totalSpent := spentPrimary + spentSecondary
	diceResult := throwExplodingDice(diceFaces, randomizer)
	total := totalSpent + diceResult

	fmt.Printf("\t\t\tResult: Spent %d + Rolled %d = %d vs Difficulty %d ... ",
		totalSpent, diceResult, total, card.Value)

	// Log details
	details := fmt.Sprintf("Card: %s (Val: %d) | Spent: %d | Rolled: %d", cardDesc, card.Value, totalSpent, diceResult)

	if total > card.Value {
		fmt.Printf("%s\n", green("SUCCESS!"))
		p.RoundScore += card.Value
		fmt.Printf("\t\t\t%s (+%d score)\n", yellow("Round Score Increased"), card.Value)
		logger.LogEvent(round, p, actionName, "Resolve", details, "SUCCESS", "")
		return true
	} else {
		fmt.Printf("%s\n", red("FAILURE."))
		p.Panic++
		fmt.Printf("\t\t\tPanic increased to %d\n", p.Panic)
		logger.LogEvent(round, p, actionName, "Resolve", details, "FAILURE", "Panic Increased")
		if p.Panic >= 3 {
			handlePanicTrigger(p, logger, round)
		}
		return false
	}
}

func printPlayerStatus(p *player) {
	fmt.Printf("\t\t\t%s\n", bold(cyan(fmt.Sprintf("--- %s STATUS (Score: %d) ---", p.Id, p.RoundScore))))
	fmt.Printf("\t\t\tO2 Cards (HP): %s | Treasure: %s | Panic: ", cyan(fmt.Sprintf("%d", len(p.O2))), bold(yellow(fmt.Sprintf("%d", p.Treasure))))
	panicColor := green
	if p.Panic == 1 {
		panicColor = yellow
	} else if p.Panic >= 2 {
		panicColor = red
	}
	fmt.Printf("%s (Dice: d%d)\n", panicColor(fmt.Sprintf("%d/3", p.Panic)), getDiceFaces(p.Panic))
	fmt.Printf("\t\t\tPools: [Enc: %d] [Env: %d] [Tec: %d] [Sop: %d]\n",
		p.AbilityPool[Encounter], p.AbilityPool[Environment], p.AbilityPool[Technical], p.AbilityPool[Soprannatural])
	fmt.Printf("\t\t\tInventory (%d/%d): ", len(p.Inventory), p.MaxInventorySize)
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
	fmt.Printf("\t\t\t%s\n", cyan("-------------------"))
}

// --- MAIN ---

func main() {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	reader := bufio.NewReader(os.Stdin)

	// INIZIALIZZA LOGGER EXCEL
	gameLogger := NewGameLogger()
	excelFileName := "game_logs.xlsx"

	// --- OGGETTI ---
	commonItemsPool := []item{
		{Type: "Coltello arrugginito", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Encounter, value: 1}}},
		{Type: "Vecchia Torcia", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Environment, value: 1}}},
		{Type: "Attrezzi base", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Technical, value: 1}}},
		{Type: "Amuleto di latta", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Soprannatural, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
		{Type: "1 Moneta", Effects: []itemEffect{{effectType: Treasure, value: 1}}},
	}
	uncommonItemsPool := []item{
		{Type: "Fiocina", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Encounter, value: 2}}},
		{Type: "Pinne Pro", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Environment, value: 2}}},
		{Type: "Attrezzi Pro", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Technical, value: 2}}},
		{Type: "Amuleto delle sirene", Effects: []itemEffect{{effectType: AddAbilityPoints, targetType: Soprannatural, value: 2}}},
		{Type: "Sacca Extra", Effects: []itemEffect{{effectType: AddSlotInventory, value: 1}}},
		{Type: "Bolla d'aria", Effects: []itemEffect{{effectType: RecoverDiscardedCards, value: 2}}},
		{Type: "3 Monete", Effects: []itemEffect{{effectType: Treasure, value: 3}}},
		{Type: "3 Monete", Effects: []itemEffect{{effectType: Treasure, value: 3}}},
	}
	rareItemsPool := []item{
		{Type: "Tridente", Effects: []itemEffect{
			{effectType: AddAbilityPoints, targetType: Encounter, value: 2},
			{effectType: AddAbilityPoints, targetType: Soprannatural, value: 2},
		}},
		{Type: "Drone Subacqueo", Effects: []itemEffect{
			{effectType: AddAbilityPoints, targetType: Technical, value: 2},
			{effectType: AddAbilityPoints, targetType: Environment, value: 2},
		}},
		{Type: "Bombola di Riserva", Effects: []itemEffect{{effectType: RecoverDiscardedCards, value: 2}}},
		{Type: "Kit Medico", Effects: []itemEffect{{effectType: ReducePanic, value: 1}}},
		{Type: "5 Monete", Effects: []itemEffect{{effectType: Treasure, value: 5}}},
		{Type: "Respiratore Pro", Effects: []itemEffect{{effectType: FreeBreath, value: 1}}},
	}
	veryRareItemsPool := []item{
		{Type: "Sonar", Effects: []itemEffect{{effectType: LookAndReorder, value: 3}}},
		{Type: "Generatore O2", Effects: []itemEffect{{effectType: RecoverDiscardedCards, value: 3}}},
		{Type: "10 Monete", Effects: []itemEffect{{effectType: Treasure, value: 10}}},
	}
	legendaryItemsPool := []item{
		{Type: "Adrenalina Pura", Effects: []itemEffect{{effectType: ReducePanic, value: 2}}},
	}

	itemsDeck := make([]item, 0)
	addItems := func(pool []item, count int) {
		for i := 0; i < count; i++ {
			itemsDeck = append(itemsDeck, pool[randomizer.Intn(len(pool))])
		}
	}
	addItems(commonItemsPool, 12)
	addItems(uncommonItemsPool, 8)
	addItems(rareItemsPool, 6)
	addItems(veryRareItemsPool, 3)
	addItems(legendaryItemsPool, 1)

	// --- COSTRUZIONE MAZZO O2 ---
	cardValuesPattern := []int{1, 1, 2, 2, 2, 3, 3, 3, 4, 5}
	o2Deck := make([]O2, 40)
	for i := range len(o2Deck) {
		o2 := O2{}
		switch {
		case i < 10:
			o2.Type = Encounter
		case i < 20:
			o2.Type = Environment
		case i < 30:
			o2.Type = Technical
		default:
			o2.Type = Soprannatural
		}
		o2.Value = cardValuesPattern[i%10]
		o2.Name = fmt.Sprintf("%s-%d", o2.Type, o2.Value)
		o2Deck[i] = o2
	}

	// --- BOSS ---
	o2Deck = append(o2Deck, O2{Name: "LEVIATHAN", Type: Encounter, SecondaryType: Soprannatural, Value: 8})
	o2Deck = append(o2Deck, O2{Name: "ABYSS STORM", Type: Environment, SecondaryType: Technical, Value: 7})
	o2Deck = append(o2Deck, O2{Name: "ALIEN WRECK", Type: Technical, SecondaryType: Encounter, Value: 7})
	o2Deck = append(o2Deck, O2{Name: "CURSE", Type: Soprannatural, SecondaryType: Environment, Value: 6})

	fmt.Println(purple("Boss Cards inserted into the deck!"))

	players := make([]player, NumberOfPlayers)
	for i := range NumberOfPlayers {
		p := player{
			Id: fmt.Sprintf("P%d", i+1),
			AbilityPool: map[abilityType]int{
				Encounter:     4,
				Environment:   4,
				Technical:     4,
				Soprannatural: 4,
			},
			Panic:            0,
			O2:               make([]O2, len(o2Deck)),
			Inventory:        make([]item, 0),
			MaxInventorySize: 3,
		}
		copy(p.O2, o2Deck)
		randomizer.Shuffle(len(p.O2), func(i, j int) { p.O2[i], p.O2[j] = p.O2[j], p.O2[i] })
		players[i] = p
	}

	for game := range NumberOfGames {
		gameLogger.InitSheet(game + 1)

		gameItems := make([]item, len(itemsDeck))
		copy(gameItems, itemsDeck)
		randomizer.Shuffle(len(gameItems), func(i, j int) { gameItems[i], gameItems[j] = gameItems[j], gameItems[i] })

		fmt.Printf("%s\n", bold(cyan(fmt.Sprintf("START GAME %d", game+1))))
		round := 1

		for playersAlive(players) > 0 && len(gameItems) > 0 {
			fmt.Printf("\t%s\n", cyan(fmt.Sprintf("--- START ROUND %d ---", round)))

			if round > 1 {
				fmt.Printf("\t%s\n", purple(">>> Resting phase..."))
				for i := range players {
					if len(players[i].O2) > 0 {
						if players[i].Panic == 0 {
							rndAbility := allAbilities[randomizer.Intn(len(allAbilities))]
							players[i].AbilityPool[rndAbility]++
							msg := fmt.Sprintf("Recovered 1 point in %s", rndAbility)
							fmt.Printf("\t\t%s is calm. %s (Total: %d)\n",
								players[i].Id, msg, players[i].AbilityPool[rndAbility])
							gameLogger.LogEvent(round, &players[i], "Rest", "Regen", string(rndAbility), "SUCCESS", msg)
						} else {
							fmt.Printf("\t\t%s is too stressed to rest! (Panic: %d)\n",
								red(players[i].Id), players[i].Panic)
							gameLogger.LogEvent(round, &players[i], "Rest", "Regen", "Panic Too High", "SKIPPED", "")
						}
						players[i].RoundScore = 0
					}
				}
			}

			for i := range players {
				if len(players[i].O2) == 0 {
					continue
				}
				p := &players[i]
				fmt.Printf("\t\t%s\n", bold(cyan(fmt.Sprintf("TURN: %s", p.Id))))

				cards, ok := drawO2CardForBreath(p)
				if !ok {
					fmt.Printf("\t\t%s\n", red(bold(fmt.Sprintf("*** %s SUFFOCATED! ***", p.Id))))
					gameLogger.LogEvent(round, p, "Turn", "Breath", "No O2 Cards", "DEATH", "Suffocated")
					continue
				}

				// Log Breath Check
				resolveCardInteraction(p, cards[0], randomizer, "BREATH CHECK", gameLogger, round)

				actionsLeft := 3
				if p.Panic >= 3 {
					actionsLeft = 0
				}

				for a := 1; a <= actionsLeft; a++ {
					printPlayerStatus(p)
					fmt.Printf("\t\t\tACTION %d/%d: (E)xplore - (U)se item - (C)alm down - (P)ass: ", a, actionsLeft)
					input, _ := reader.ReadString('\n')
					input = strings.TrimSpace(strings.ToLower(input))

					switch input {
					case "c":
						if p.Panic > 0 {
							targetDifficulty := 6
							willpowerDice := 6
							fmt.Printf("\t\t\t%s (Target: %d)\n", bold(cyan("CALM DOWN CHECK")), targetDifficulty)
							totalSpent := 0
							spentMap := make(map[abilityType]int)
							fmt.Println("\t\t\tChoose points to spend from each pool:")
							order := []abilityType{Encounter, Environment, Technical, Soprannatural}
							for _, abType := range order {
								available := p.AbilityPool[abType]
								if available == 0 {
									continue
								}
								validInput := false
								for !validInput {
									fmt.Printf("\t\t\t- %s (Available: %d): ", abType, available)
									input, _ := reader.ReadString('\n')
									input = strings.TrimSpace(input)
									if input == "" {
										input = "0"
									}
									val, err := strconv.Atoi(input)
									if err == nil && val >= 0 && val <= available {
										spentMap[abType] = val
										totalSpent += val
										validInput = true
									} else {
										fmt.Printf("\t\t\t  Invalid amount (0-%d).\n", available)
									}
								}
							}
							for abType, amount := range spentMap {
								p.AbilityPool[abType] -= amount
							}
							roll := randomizer.Intn(willpowerDice) + 1
							total := totalSpent + roll

							details := fmt.Sprintf("Spent: %d | Rolled: %d (d%d)", totalSpent, roll, willpowerDice)

							fmt.Printf("\t\t\tResult: %s = %d vs Target %d ... ", details, total, targetDifficulty)
							if total >= targetDifficulty {
								p.Panic--
								fmt.Printf("%s (Panic is now %d)\n", green("SUCCESS!"), p.Panic)
								gameLogger.LogEvent(round, p, "Action", "Calm Down", details, "SUCCESS", "Panic Reduced")
							} else {
								fmt.Printf("%s (Panic remains %d)\n", red("FAILED."), p.Panic)
								gameLogger.LogEvent(round, p, "Action", "Calm Down", details, "FAILURE", "Panic Unchanged")
							}
						} else {
							fmt.Println("\t\t\tYou are already calm.")
							a--
						}

					case "u":
						if len(p.Inventory) == 0 {
							fmt.Println("\t\t\tInventory empty.")
							a--
							continue
						}
						fmt.Println("\t\t\tSelect item to use:")
						for idx, it := range p.Inventory {
							desc := it.Type
							if len(it.Effects) > 0 {
								desc += fmt.Sprintf(" [%s %d]", it.Effects[0].effectType, it.Effects[0].value)
							}
							fmt.Printf("\t\t\t[%d] %s\n", idx+1, cyan(desc))
						}
						fmt.Printf("\t\t\tChoose [1-%d] (or 0 to cancel): ", len(p.Inventory))
						input, _ := reader.ReadString('\n')
						input = strings.TrimSpace(input)
						choice, err := strconv.Atoi(input)
						if err == nil && choice > 0 && choice <= len(p.Inventory) {
							index := choice - 1
							item := p.Inventory[index]

							gameLogger.LogEvent(round, p, "Action", "Use Item", item.Type, "USED", "")

							fmt.Printf("\t\t\tUsing %s...\n", item.Type)
							for _, eff := range item.Effects {
								switch eff.effectType {
								case AddAbilityPoints:
									p.AbilityPool[eff.targetType] += eff.value
									fmt.Printf("\t\t\tRestored %d points to %s\n", eff.value, eff.targetType)
								case ReducePanic:
									if p.Panic > 0 {
										p.Panic -= eff.value
										if p.Panic < 0 {
											p.Panic = 0
										}
										fmt.Printf("\t\t\tPanic reduced by %d\n", eff.value)
									} else {
										fmt.Println("\t\t\t(Panic was already 0)")
									}
								case AddSlotInventory:
									p.MaxInventorySize += eff.value
									fmt.Println("\t\t\tInventory expanded.")
								case RecoverDiscardedCards:
									if len(p.DiscardedO2) > 0 {
										recovered, remainingDiscarded := draw(eff.value, p.DiscardedO2)
										p.DiscardedO2 = remainingDiscarded
										p.O2 = append(p.O2, recovered...)
										fmt.Printf("\t\t\tRecovered %d O2 cards from discard pile! (HP: %d)\n", len(recovered), len(p.O2))
									} else {
										fmt.Println("\t\t\tNo cards to recover in discard pile.")
									}
								}
							}
							p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
						} else if choice == 0 {
							fmt.Println("\t\t\tCancelled.")
							a--
						} else {
							fmt.Println("\t\t\tInvalid selection.")
							a--
						}

					case "e":
						cards, newO2 := draw(1, p.O2)
						if len(cards) > 0 {
							p.O2 = newO2
							p.DiscardedO2 = append(p.DiscardedO2, cards...)
							resolveCardInteraction(p, cards[0], randomizer, "EXPLORE", gameLogger, round)
						} else {
							fmt.Println("\t\t\tNo more O2 cards to explore.")
						}
					case "p":
						actionsLeft = 0
						gameLogger.LogEvent(round, p, "Action", "Pass", "", "PASSED", "")
					default:
						fmt.Println("\t\t\tInvalid command.")
						a--
					}
				}
				p.Effects = make([]itemEffect, 0)
			}

			fmt.Printf("\t\t%s\n", bold(yellow("--- LOOT PHASE (DRAFT) ---")))

			type playerRef struct {
				Index int
				Score int
				Id    string
			}
			var ranking []playerRef
			for i, p := range players {
				if len(p.O2) > 0 && p.RoundScore > 0 {
					ranking = append(ranking, playerRef{Index: i, Score: p.RoundScore, Id: p.Id})
				}
			}
			sort.Slice(ranking, func(i, j int) bool {
				if ranking[i].Score == ranking[j].Score {
					return players[ranking[i].Index].Panic < players[ranking[j].Index].Panic
				}
				return ranking[i].Score > ranking[j].Score
			})

			if len(ranking) == 0 {
				fmt.Println("\t\tNo items found (No one scored points).")
				gameLogger.LogEvent(round, nil, "Loot", "None", "No Score", "", "")
			} else {
				itemsToDraw := len(ranking)
				if itemsToDraw == 1 {
					itemsToDraw = 2
				}
				marketItems, remainingDeck := draw(itemsToDraw, gameItems)
				gameItems = remainingDeck
				fmt.Printf("\t\tItems found: %d\n", len(marketItems))

				for i, rankRef := range ranking {
					if len(marketItems) == 0 {
						break
					}
					p := &players[rankRef.Index]
					fmt.Printf("\t\t%s (Score: %d) is choosing...\n", cyan(p.Id), rankRef.Score)
					choiceIndex := 0
					if len(marketItems) > 1 {
						validChoice := false
						for !validChoice {
							fmt.Println("\t\tAvailable Items:")
							for idx, it := range marketItems {
								desc := it.Type
								if len(it.Effects) > 0 {
									desc += fmt.Sprintf(" (%s %d)", it.Effects[0].effectType, it.Effects[0].value)
								}
								fmt.Printf("\t\t\t[%d] %s\n", idx+1, yellow(desc))
							}
							fmt.Printf("\t\tChoose item [1-%d]: ", len(marketItems))
							input, _ := reader.ReadString('\n')
							input = strings.TrimSpace(input)
							val, err := strconv.Atoi(input)
							if err == nil && val >= 1 && val <= len(marketItems) {
								choiceIndex = val - 1
								validChoice = true
							} else {
								fmt.Println("\t\tInvalid selection.")
							}
						}
					} else {
						fmt.Printf("\t\tOnly one item left: %s. Auto-looting.\n", marketItems[0].Type)
						choiceIndex = 0
					}
					selectedItem := marketItems[choiceIndex]
					marketItems = append(marketItems[:choiceIndex], marketItems[choiceIndex+1:]...)

					isTreasure := false
					treasureVal := 0
					for _, ef := range selectedItem.Effects {
						if ef.effectType == Treasure {
							isTreasure = true
							treasureVal += ef.value
						}
					}

					logDetails := selectedItem.Type

					if isTreasure {
						p.Treasure += treasureVal
						fmt.Printf("\t\t%s obtained Treasure: %s (+%d)\n", p.Id, yellow(selectedItem.Type), treasureVal)
						logDetails += fmt.Sprintf(" (Treasure +%d)", treasureVal)
					} else {
						if len(p.Inventory) < p.MaxInventorySize {
							p.Inventory = append(p.Inventory, selectedItem)
							fmt.Printf("\t\t%s obtained Item: %s\n", p.Id, cyan(selectedItem.Type))
						} else {
							fmt.Printf("\t\tInventory FULL. Swap with an existing item? (y/n): ")
							input, _ := reader.ReadString('\n')
							if strings.TrimSpace(strings.ToLower(input)) == "y" {
								for idx, invItem := range p.Inventory {
									fmt.Printf("\t\t\t[%d] %s\n", idx+1, invItem.Type)
								}
								fmt.Print("\t\tDrop which item? [1-N]: ")
								input, _ = reader.ReadString('\n')
								dropIdx, err := strconv.Atoi(strings.TrimSpace(input))
								if err == nil && dropIdx >= 1 && dropIdx <= len(p.Inventory) {
									dropped := p.Inventory[dropIdx-1]
									p.Inventory[dropIdx-1] = selectedItem
									fmt.Printf("\t\tDropped %s and took %s.\n", dropped.Type, selectedItem.Type)
									logDetails += fmt.Sprintf(" (Swapped with %s)", dropped.Type)
								} else {
									fmt.Println("\t\tInvalid input. Item discarded.")
									logDetails += " (Discarded Full Inv)"
								}
							} else {
								fmt.Printf("\t\t%s discarded due to full inventory.\n", selectedItem.Type)
								logDetails += " (Discarded Full Inv)"
							}
						}
					}
					gameLogger.LogEvent(round, p, "Loot", "Acquire", logDetails, "SUCCESS", "")

					if len(ranking) == 1 && i == 0 {
						fmt.Println("\t\tRemaining items are discarded.")
						marketItems = make([]item, 0)
					}
				}
			}
			round++
		}

		fmt.Printf("%s\n", bold(cyan(fmt.Sprintf("GAME OVER"))))
		var survivors []player
		for _, p := range players {
			if len(p.O2) > 0 {
				survivors = append(survivors, p)
			}
		}
		if len(survivors) == 0 {
			fmt.Println(red("\tAll players perished in the deep. The ocean claims all."))
			gameLogger.LogEvent(round, nil, "GameOver", "End", "No Survivors", "LOSS", "")
		} else {
			sort.Slice(survivors, func(i, j int) bool {
				if survivors[i].Treasure == survivors[j].Treasure {
					return len(survivors[i].O2) > len(survivors[j].O2)
				}
				return survivors[i].Treasure > survivors[j].Treasure
			})
			winner := survivors[0]
			fmt.Println()
			fmt.Printf("\t%s\n", bold(green("🏆 WE HAVE A WINNER! 🏆")))
			fmt.Printf("\t%s survived with %s Coins!\n", bold(winner.Id), bold(yellow(fmt.Sprintf("%d", winner.Treasure))))

			msg := fmt.Sprintf("Winner: %s (Treasure: %d)", winner.Id, winner.Treasure)
			gameLogger.LogEvent(round, &winner, "GameOver", "Win", msg, "VICTORY", "")

			if len(survivors) > 1 {
				fmt.Println("\tOther survivors:")
				for i := 1; i < len(survivors); i++ {
					p := survivors[i]
					fmt.Printf("\t%d. %s (%d Coins)\n", i+1, p.Id, p.Treasure)
				}
			}
		}
	}

	// Salvataggio finale Excel
	gameLogger.Save(excelFileName)
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

func draw[T any](n int, slice []T) ([]T, []T) {
	if len(slice) <= n {
		n = len(slice)
	}
	return slice[:n], slice[n:]
}
