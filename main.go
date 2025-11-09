package main

import "board-game-course/model"

var NumberOfPlayers = model.NewGameParameter(model.NumberOfPlayers, 2)
var NumberOfPanicCardsWithOneType = model.NewGameParameter(model.NumberOfPanicCardsWithOneType, 12)
var NumberOfPanicCardsWithTwoType = model.NewGameParameter(model.NumberOfPanicCardsWithTwoType, 15)
var NumberOfPanicCardsWithThreeType = model.NewGameParameter(model.NumberOfPanicCardsWithThreeType, 2)
var NumberOfPanicCardsToActivateEffect = model.NewGameParameter(model.NumberOfPanicCardsToActivateEffect, 3)
var NumberOfItemCardsForCommonRarity = model.NewGameParameter(model.NumberOfItemCardsForCommonRarity, 8)
var NumberOfItemCardsForUncommonRarity = model.NewGameParameter(model.NumberOfItemCardsForUncommonRarity, 6)
var NumberOfItemCardsForRareRarity = model.NewGameParameter(model.NumberOfItemCardsForRareRarity, 4)
var NumberOfItemCardsForLegendaryRarity = model.NewGameParameter(model.NumberOfItemCardsForLegendaryRarity, 2)
var NumberOfItemSlots = model.NewGameParameter(model.NumberOfItemSlots, 3)
var NumberOfAmulets = model.NewGameParameter(model.NumberOfAmulets, 4)
var NumberOfAmuletsToWin = model.NewGameParameter(model.NumberOfAmuletsToWin, 3)

var NumberOfGames = 1

func main() {

	model.NewGame(
		NumberOfPlayers,
		NumberOfPanicCardsWithOneType,
		NumberOfPanicCardsWithTwoType,
		NumberOfPanicCardsWithThreeType,
		NumberOfPanicCardsToActivateEffect,
		NumberOfItemCardsForCommonRarity,
		NumberOfItemCardsForUncommonRarity,
		NumberOfItemCardsForRareRarity,
		NumberOfItemCardsForLegendaryRarity,
		NumberOfItemSlots,
		NumberOfAmulets,
		NumberOfAmuletsToWin,
	)

	//game.Run(NumberOfGames)
}
