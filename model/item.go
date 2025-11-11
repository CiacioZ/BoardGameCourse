package model

type itemCardRarity string

const (
	Common    itemCardRarity = "COMMON"
	Uncommon  itemCardRarity = "UNCOMMON"
	Rare      itemCardRarity = "RARE"
	Legendary itemCardRarity = "LEGENDARY"
)

type itemCard struct {
	genericCard
	rarity itemCardRarity
	items  map[int]item
}

type itemEffectType string

const (
	LookNextO2Cards         itemEffectType = "LOOK_NEXT_O2_CARDS"
	MovementCostReduction   itemEffectType = "MOVEMENT_COST_REDUCTION"
	BreathCostReduction     itemEffectType = "BREATH_COST_REDUCTION"
	BlockPlayer             itemEffectType = "BLOCK_PLAYER"
	IgnorePanicActivation   itemEffectType = "IGNORE_PANIC_ACTIVATION"
	AnotherPlayerMustDrawO2 itemEffectType = "ANOTHER_PLAYER_MUST_DRAW_O2"
	StealItemFromPlayer     itemEffectType = "STEAL_ITEM_FROM_PLAYER"
	StealAmuletFromPLayer   itemEffectType = "STEAL_AMULET_FROM_PLAYER"
	RecoverDiscardedO2      itemEffectType = "RECOVER_DISCARDED_O2"
	ReorderNextO2Cards      itemEffectType = "REORDER_NEXT_O2_CARDS"
)

type itemEffect struct {
	effectType itemEffectType
	value      int
}

type itemType string

const (
	Utility       itemType = "UTILITY"
	TreasureToken itemType = "TREASURE_TOKEN"
	Amulets       itemType = "AMULETS"
)

type item struct {
	name     string
	itemType itemType
	effects  []itemEffect
	quantity int
}

var flashlight = item{
	name:     "flashlight",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: LookNextO2Cards,
			value:      2,
		},
	},
}

var enhancedFins = item{
	name:     "EnhancedFins",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: MovementCostReduction,
			value:      1,
		},
	},
}

var advancedMask = item{
	name:     "AdvancedMask",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: BreathCostReduction,
			value:      1,
		},
	},
}

var net = item{
	name:     "Net",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: BlockPlayer,
			value:      1,
		},
	},
}

var reinforcedNet = item{
	name:     "ReinforcedNet",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: BlockPlayer,
			value:      2,
		},
	},
}

var antistressKit = item{
	name:     "AntistressKit",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: IgnorePanicActivation,
			value:      1,
		},
	},
}

var spearGun = item{
	name:     "SpearGun",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: AnotherPlayerMustDrawO2,
			value:      2,
		},
	},
}

var harpoon = item{
	name:     "Harpoon",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: StealItemFromPlayer,
			value:      1,
		},
	},
}

var mysticHarpoon = item{
	name:     "MysticHarpoon",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: StealAmuletFromPLayer,
			value:      1,
		},
	},
}

var emergencyAirBag = item{
	name:     "EmergencyAirBag",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: RecoverDiscardedO2,
			value:      3,
		},
	},
}

var sonar = item{
	name:     "Sonar",
	itemType: Utility,
	quantity: 1,
	effects: []itemEffect{
		{
			effectType: ReorderNextO2Cards,
			value:      3,
		},
	},
}

var amulet = item{
	name:     "Amulet",
	itemType: Amulets,
	quantity: 1,
	effects:  []itemEffect{},
}

var smallTreasure = item{
	name:     "Treasure",
	itemType: TreasureToken,
	effects:  []itemEffect{},
	quantity: 1,
}

var mediumTreasure = item{
	name:     "Treasure",
	itemType: TreasureToken,
	effects:  []itemEffect{},
	quantity: 3,
}

var bigTreasure = item{
	name:     "Treasure",
	itemType: TreasureToken,
	effects:  []itemEffect{},
	quantity: 5,
}

var commonItemCards = []itemCard{
	{
		genericCard: genericCard{
			Name: "Common Card 1",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: flashlight,
			2: flashlight,
			3: flashlight,
			4: net,
			5: net,
			6: net,
			7: antistressKit,
			8: antistressKit,
			9: antistressKit,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 2",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: enhancedFins,
			2: enhancedFins,
			3: enhancedFins,
			4: advancedMask,
			5: advancedMask,
			6: advancedMask,
			7: reinforcedNet,
			8: reinforcedNet,
			9: reinforcedNet,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 3",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: smallTreasure,
			2: smallTreasure,
			3: smallTreasure,
			4: mediumTreasure,
			5: mediumTreasure,
			6: mediumTreasure,
			7: bigTreasure,
			8: bigTreasure,
			9: bigTreasure,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 4",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: smallTreasure,
			2: smallTreasure,
			3: smallTreasure,
			4: mediumTreasure,
			5: mediumTreasure,
			6: mediumTreasure,
			7: bigTreasure,
			8: bigTreasure,
			9: bigTreasure,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 5",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: flashlight,
			2: flashlight,
			3: flashlight,
			4: net,
			5: net,
			6: net,
			7: antistressKit,
			8: antistressKit,
			9: antistressKit,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 6",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: enhancedFins,
			2: enhancedFins,
			3: enhancedFins,
			4: advancedMask,
			5: advancedMask,
			6: advancedMask,
			7: reinforcedNet,
			8: reinforcedNet,
			9: reinforcedNet,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 7",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: advancedMask,
			2: advancedMask,
			3: advancedMask,
			4: reinforcedNet,
			5: reinforcedNet,
			6: reinforcedNet,
			7: enhancedFins,
			8: enhancedFins,
			9: enhancedFins,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 8",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: net,
			2: net,
			3: net,
			4: antistressKit,
			5: antistressKit,
			6: antistressKit,
			7: flashlight,
			8: flashlight,
			9: flashlight,
		},
	},
}

var uncommonItemCards = []itemCard{
	{
		genericCard: genericCard{
			Name: "Uncommon Card 1",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: smallTreasure,
			2: smallTreasure,
			3: smallTreasure,
			4: mediumTreasure,
			5: mediumTreasure,
			6: mediumTreasure,
			7: bigTreasure,
			8: bigTreasure,
			9: bigTreasure,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 2",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: reinforcedNet,
			2: reinforcedNet,
			3: reinforcedNet,
			4: enhancedFins,
			5: enhancedFins,
			6: enhancedFins,
			7: advancedMask,
			8: advancedMask,
			9: advancedMask,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 3",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: antistressKit,
			2: antistressKit,
			3: antistressKit,
			4: flashlight,
			5: flashlight,
			6: flashlight,
			7: net,
			8: net,
			9: net,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 4",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: smallTreasure,
			2: smallTreasure,
			3: smallTreasure,
			4: mediumTreasure,
			5: mediumTreasure,
			6: mediumTreasure,
			7: bigTreasure,
			8: bigTreasure,
			9: bigTreasure,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 5",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: reinforcedNet,
			2: reinforcedNet,
			3: reinforcedNet,
			4: enhancedFins,
			5: enhancedFins,
			6: enhancedFins,
			7: advancedMask,
			8: advancedMask,
			9: advancedMask,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 6",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: antistressKit,
			2: antistressKit,
			3: antistressKit,
			4: flashlight,
			5: flashlight,
			6: flashlight,
			7: net,
			8: net,
			9: net,
		},
	},
}

var rareItemCards = []itemCard{
	{
		genericCard: genericCard{
			Name: "Rare Card 1",
			Type: ItemType,
		},
		rarity: Rare,
		items: map[int]item{
			1: mediumTreasure,
			2: mediumTreasure,
			3: mediumTreasure,
			4: bigTreasure,
			5: bigTreasure,
			6: bigTreasure,
			7: amulet,
			8: amulet,
			9: amulet,
		},
	},
	{
		genericCard: genericCard{
			Name: "Rare Card 2",
			Type: ItemType,
		},
		rarity: Rare,
		items: map[int]item{
			1: harpoon,
			2: harpoon,
			3: harpoon,
			4: spearGun,
			5: spearGun,
			6: spearGun,
			7: amulet,
			8: amulet,
			9: amulet,
		},
	},
	{
		genericCard: genericCard{
			Name: "Rare Card 3",
			Type: ItemType,
		},
		rarity: Rare,
		items: map[int]item{
			1: mediumTreasure,
			2: mediumTreasure,
			3: mediumTreasure,
			4: amulet,
			5: amulet,
			6: amulet,
			7: bigTreasure,
			8: bigTreasure,
			9: bigTreasure,
		},
	},
	{
		genericCard: genericCard{
			Name: "Rare Card 4",
			Type: ItemType,
		},
		rarity: Rare,
		items: map[int]item{
			1: spearGun,
			2: spearGun,
			3: spearGun,
			4: harpoon,
			5: harpoon,
			6: harpoon,
			7: amulet,
			8: amulet,
			9: amulet,
		},
	},
}

var legendaryItemCards = []itemCard{
	{
		genericCard: genericCard{
			Name: "Legendary Card 1",
			Type: ItemType,
		},
		rarity: Legendary,
		items: map[int]item{
			1: bigTreasure,
			2: bigTreasure,
			3: bigTreasure,
			4: sonar,
			5: sonar,
			6: sonar,
			7: amulet,
			8: amulet,
			9: amulet,
		},
	},
	{
		genericCard: genericCard{
			Name: "Legendary Card 2",
			Type: ItemType,
		},
		rarity: Legendary,
		items: map[int]item{
			1: emergencyAirBag,
			2: emergencyAirBag,
			3: emergencyAirBag,
			4: mysticHarpoon,
			5: mysticHarpoon,
			6: mysticHarpoon,
			7: amulet,
			8: amulet,
			9: amulet,
		},
	},
}
