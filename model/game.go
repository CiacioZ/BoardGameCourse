package model

import "fmt"

type GameParameterType string

const (
	NumberOfPlayers                     GameParameterType = "NUMBER_OF_PLAYERS"
	NumberOfPanicCardsWithOneType       GameParameterType = "NUMBER_OF_PANIC_CARD_WITH_ONE_TYPE"
	NumberOfPanicCardsWithTwoType       GameParameterType = "NUMBER_OF_PANIC_CARD_WITH_TWO_TYPE"
	NumberOfPanicCardsWithThreeType     GameParameterType = "NUMBER_OF_PANIC_CARD_WITH_THREE_TYPE"
	NumberOfPanicCardsToActivateEffect  GameParameterType = "NUMBER_OF_PANIC_CARD_TO_ACTIVATE_EFFECT"
	NumberOfItemCardsForCommonRarity    GameParameterType = "NUMBER_OF_ITEM_CARD_FOR_COMMON_RARITY"
	NumberOfItemCardsForUncommonRarity  GameParameterType = "NUMBER_OF_ITEM_CARD_FOR_UNCOMMON_RARITY"
	NumberOfItemCardsForRareRarity      GameParameterType = "NUMBER_OF_ITEM_CARD_FOR_RARE_RARITY"
	NumberOfItemCardsForLegendaryRarity GameParameterType = "NUMBER_OF_ITEM_CARD_FOR_LEGENDARY_RARITY"
	NumberOfItemSlots                   GameParameterType = "NUMBER_OF_ITEM_SLOTS"
	NumberOfAmulets                     GameParameterType = "NUMBER_OF_AMULETS"
	NumberOfAmuletsToWin                GameParameterType = "NUMBER_OF_AMULETS_TO_WIN"
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
	}

	for i := range g.parameters.values[NumberOfPlayers] {
		player := NewPlayer(fmt.Sprintf("Player %d", i))

		oxygenDeck := g.GenerateOxygenDeck()

		player.OxygenCards = oxygenDeck

		g.state.AddPlayer(player)
	}

	return g
}

func (g *game) Run(numberOfGames int) {

	g.state.Round = 1

	for !g.IsGameEnded() {
		fmt.Println(fmt.Printf("Start Round: %d", g.state.Round))

		//TODO: Insert logic

		fmt.Println(fmt.Printf("End Round: %d", g.state.Round))

		g.state.Round++
	}

}

func (g *game) IsDavyJonesIsDead() bool {
	return false
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
	for i := 1; i <= g.parameters.values[NumberOfPanicCardsWithOneType]; i++ {

		deck = append(deck, panicCard{
			genericCard: genericCard{
				Name: fmt.Sprintf("Panic Card 1C %d", i),
				Type: PanicType,
			},
		})
	}

	//Generate two color type panic card
	for i := 1; i <= g.parameters.values[NumberOfPanicCardsWithTwoType]; i++ {

	}

	//Generate three color type panic card
	for i := 1; i <= g.parameters.values[NumberOfPanicCardsWithThreeType]; i++ {

	}

	//Generate common item card
	for i := 1; i <= g.parameters.values[NumberOfItemCardsForCommonRarity]; i++ {

	}

	//Generate uncommon item card
	for i := 1; i <= g.parameters.values[NumberOfItemCardsForUncommonRarity]; i++ {

	}

	//Generate rare item card
	for i := 1; i <= g.parameters.values[NumberOfItemCardsForRareRarity]; i++ {

	}

	//Generate legendary item card
	for i := 1; i <= g.parameters.values[NumberOfItemCardsForLegendaryRarity]; i++ {

	}

	return deck
}
