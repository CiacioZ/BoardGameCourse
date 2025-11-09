package main

import "board-game-course/model"

var NumberOfPlayers = model.NewGameParameter(model.NumberOfPlayers, 2)
var NumberOfPanicCardsToActivateEffect = model.NewGameParameter(model.NumberOfPanicCardsToActivateEffect, 3)
var NumberOfItemSlots = model.NewGameParameter(model.NumberOfItemSlots, 3)
var NumberOfAmuletsToWin = model.NewGameParameter(model.NumberOfAmuletsToWin, 3)

var NumberOfGames = 1

func main() {

	model.NewGame(
		NumberOfPlayers,
		NumberOfPanicCardsToActivateEffect,
		NumberOfItemSlots,
		NumberOfAmuletsToWin,
	)

	//game.Run(NumberOfGames)
}
