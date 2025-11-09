package model

type panicType string

const (
	Blue   panicType = "BLUE"
	Green  panicType = "GREEN"
	Red    panicType = "RED"
	Yellow panicType = "YELLOW"
	Black  panicType = "BLACK"
	Purple panicType = "PURPLE"
)

type panicCard struct {
	genericCard
	panicTypes []panicType
}

type panicEffectType string

const (
	MoveUp                    panicEffectType = "MOVE_UP"
	MoveDown                  panicEffectType = "MOVE_DOWN"
	CannotExplore             panicEffectType = "CANNOT_EXPLORE"
	JumpTurn                  panicEffectType = "JUMP_TURN"
	DiscardO2                 panicEffectType = "DISCARD_O2"
	DropObject                panicEffectType = "DROP_OBJECT"
	MustCalmDonw              panicEffectType = "MUST_CALM_DOWN"
	ExchangeRandomCard        panicEffectType = "EXCHANGE_RANDOM_CARD"
	DropTreasureToken         panicEffectType = "DROP_TREASURE_TOKEN"
	DropAmulet                panicEffectType = "DROP_AMULET"
	DropEverything            panicEffectType = "DROP_EVERYTHING"
	DropEverythingButAmulets  panicEffectType = "DROP_EVERYTHING_BUT_AMULETS"
	DropO2ForSameLevelPlayers panicEffectType = "DROP_O2_SAME_LEVEL_PLAYERS"
	DrawO2                    panicEffectType = "DRAW_O2"
)

type panicEffect struct {
	effectType panicEffectType
	value      int
}

var blueLowLevels = []panicEffect{
	{
		effectType: CannotExplore,
		value:      1,
	},
}
var blueMidLevels = []panicEffect{
	{
		effectType: MoveUp,
		value:      1,
	},
}
var blueDeepLevels = []panicEffect{
	{
		effectType: JumpTurn,
		value:      1,
	},
}

var redLowLevels = []panicEffect{
	{
		effectType: DiscardO2,
		value:      2,
	},
}
var redMidLevels = []panicEffect{
	{
		effectType: DropObject,
		value:      1,
	},
	{
		effectType: MoveUp,
		value:      2,
	},
	{
		effectType: MustCalmDonw,
		value:      1,
	},
}
var redDeepLevels = []panicEffect{
	{
		effectType: DropObject,
		value:      1,
	},
	{
		effectType: MoveUp,
		value:      3,
	},
	{
		effectType: MustCalmDonw,
		value:      1,
	},
	{
		effectType: DiscardO2,
		value:      2,
	},
}

var greenLowLevels = []panicEffect{
	{
		effectType: MoveUp,
		value:      2,
	},
}
var greenMidLevels = []panicEffect{
	{
		effectType: MoveUp,
		value:      2,
	},
	{
		effectType: DiscardO2,
		value:      1,
	},
	{
		effectType: DropO2ForSameLevelPlayers,
		value:      1,
	},
}
var greenDeepLevels = []panicEffect{
	{
		effectType: MoveUp,
		value:      3,
	},
	{
		effectType: DiscardO2,
		value:      2,
	},
	{
		effectType: DropO2ForSameLevelPlayers,
		value:      2,
	},
}

var yellowLowLevels = []panicEffect{
	{
		effectType: DrawO2,
		value:      2,
	},
}
var yellowMidLevels = []panicEffect{
	{
		effectType: DrawO2,
		value:      3,
	},
	{
		effectType: MoveUp,
		value:      1,
	},
}
var yellowDeepLevels = []panicEffect{
	{
		effectType: DrawO2,
		value:      5,
	},
	{
		effectType: MoveUp,
		value:      2,
	},
}

var blackLowLevels = []panicEffect{
	{
		effectType: DropObject,
		value:      1,
	},
}
var blackMidLevels = []panicEffect{
	{
		effectType: DropEverythingButAmulets,
		value:      0,
	},
}
var blackDeepLevels = []panicEffect{
	{
		effectType: DropEverything,
		value:      0,
	},
}

var purpleLowLevels = []panicEffect{
	{
		effectType: ExchangeRandomCard,
		value:      1,
	},
}
var purpleMidLevels = []panicEffect{
	{
		effectType: ExchangeRandomCard,
		value:      1,
	},
	{
		effectType: DropTreasureToken,
		value:      1,
	},
}
var purpleDeepLevels = []panicEffect{
	{
		effectType: ExchangeRandomCard,
		value:      1,
	},
	{
		effectType: DropAmulet,
		value:      1,
	},
}

var panicActivationEffects = map[panicType]map[int][]panicEffect{
	Blue: {
		1: blueLowLevels,
		2: blueLowLevels,
		3: blueLowLevels,
		4: blueMidLevels,
		5: blueMidLevels,
		6: blueMidLevels,
		7: blueDeepLevels,
		8: blueDeepLevels,
		9: blueDeepLevels,
	},
	Red: {
		1: redLowLevels,
		2: redLowLevels,
		3: redLowLevels,
		4: redMidLevels,
		5: redMidLevels,
		6: redMidLevels,
		7: redDeepLevels,
		8: redDeepLevels,
		9: redDeepLevels,
	},
	Yellow: {
		1: yellowLowLevels,
		2: yellowLowLevels,
		3: yellowLowLevels,
		4: yellowMidLevels,
		5: yellowMidLevels,
		6: yellowMidLevels,
		7: yellowDeepLevels,
		8: yellowDeepLevels,
		9: yellowDeepLevels,
	},
	Purple: {
		1: purpleLowLevels,
		2: purpleLowLevels,
		3: purpleLowLevels,
		4: purpleMidLevels,
		5: purpleMidLevels,
		6: purpleMidLevels,
		7: purpleDeepLevels,
		8: purpleDeepLevels,
		9: purpleDeepLevels,
	},
	Green: {
		1: greenLowLevels,
		2: greenLowLevels,
		3: greenLowLevels,
		4: greenMidLevels,
		5: greenMidLevels,
		6: greenMidLevels,
		7: greenDeepLevels,
		8: greenDeepLevels,
		9: greenDeepLevels,
	},
	Black: {
		1: blackLowLevels,
		2: blackLowLevels,
		3: blackLowLevels,
		4: blackMidLevels,
		5: blackMidLevels,
		6: blackMidLevels,
		7: blackDeepLevels,
		8: blackDeepLevels,
		9: blackDeepLevels,
	},
}

var singlePanicCards = []panicCard{
	{
		genericCard: genericCard{
			Name: "Blue",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Blue,
		},
	},
	{
		genericCard: genericCard{
			Name: "Red",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red,
		},
	},
	{
		genericCard: genericCard{
			Name: "Green",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Green,
		},
	},
	{
		genericCard: genericCard{
			Name: "Yellow",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Yellow,
		},
	},
	{
		genericCard: genericCard{
			Name: "Purple",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Purple,
		},
	},
	{
		genericCard: genericCard{
			Name: "Black",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Black,
		},
	},
}

/*



green – blue

green – yellow

green – black

green – purple

blue – yellow

blue – black

blue – purple

yellow – black

yellow – purple

black – purple
*/

var doublePanicCards = []panicCard{
	{
		genericCard: genericCard{
			Name: "Red - Green",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Green,
		},
	},
	{
		genericCard: genericCard{
			Name: "Red - Blue",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Blue,
		},
	},
	{
		genericCard: genericCard{
			Name: "Red - Yellow",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Yellow,
		},
	},
	{
		genericCard: genericCard{
			Name: "Red - Black",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Black,
		},
	},
	{
		genericCard: genericCard{
			Name: "Red - Purple",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Purple,
		},
	},
	{
		genericCard: genericCard{
			Name: "Green - Blue",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Green, Blue,
		},
	},
	{
		genericCard: genericCard{
			Name: "Green - Yellow",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Green, Yellow,
		},
	},
	{
		genericCard: genericCard{
			Name: "Green - Black",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Green, Black,
		},
	},
	{
		genericCard: genericCard{
			Name: "Green - Purple",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Green, Purple,
		},
	},
	{
		genericCard: genericCard{
			Name: "Green - Purple",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Purple,
		},
	},
	{
		genericCard: genericCard{
			Name: "Blue - Yellow",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Blue, Yellow,
		},
	},
	{
		genericCard: genericCard{
			Name: "Blue - Black",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Blue, Black,
		},
	},
	{
		genericCard: genericCard{
			Name: "Blue - Purple",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Blue, Purple,
		},
	},
	{
		genericCard: genericCard{
			Name: "Yellow - Black",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Yellow, Purple,
		},
	},
	{
		genericCard: genericCard{
			Name: "Yellow - Purple",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Yellow, Purple,
		},
	},
	{
		genericCard: genericCard{
			Name: "Black - Purple",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Black, Purple,
		},
	},
}

var triplePanicCards = []panicCard{
	{
		genericCard: genericCard{
			Name: "Red - Purple - Yellow",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Purple, Yellow,
		},
	},
	{
		genericCard: genericCard{
			Name: "Green - Blue - Black",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Green, Blue, Black,
		},
	},
	{
		genericCard: genericCard{
			Name: "Red - Purple - Green",
			Type: PanicType,
		},
		panicTypes: []panicType{
			Red, Purple, Green,
		},
	},
}
