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

type item struct {
	name    string
	effects []itemEffect
}

var FlashLight = item{
	name: "Flashlight",
	effects: []itemEffect{
		{
			effectType: LookNextO2Cards,
			value:      2,
		},
	},
}

var EnhancedFins = item{
	name: "EnhancedFins",
	effects: []itemEffect{
		{
			effectType: MovementCostReduction,
			value:      1,
		},
	},
}

var AdvancedMask = item{
	name: "AdvancedMask",
	effects: []itemEffect{
		{
			effectType: BreathCostReduction,
			value:      1,
		},
	},
}

var Net = item{
	name: "Net",
	effects: []itemEffect{
		{
			effectType: BlockPlayer,
			value:      1,
		},
	},
}

var ReinforcedNet = item{
	name: "ReinforcedNet",
	effects: []itemEffect{
		{
			effectType: BlockPlayer,
			value:      2,
		},
	},
}

var AntistressKit = item{
	name: "AntistressKit",
	effects: []itemEffect{
		{
			effectType: IgnorePanicActivation,
			value:      1,
		},
	},
}

var SpearGun = item{
	name: "SpearGun",
	effects: []itemEffect{
		{
			effectType: AnotherPlayerMustDrawO2,
			value:      2,
		},
	},
}

var Harpoon = item{
	name: "Harpoon",
	effects: []itemEffect{
		{
			effectType: StealItemFromPlayer,
			value:      1,
		},
	},
}

var MysticHarpoon = item{
	name: "MysticHarpoon",
	effects: []itemEffect{
		{
			effectType: StealAmuletFromPLayer,
			value:      1,
		},
	},
}

var EmergencyAirBag = item{
	name: "EmergencyAirBag",
	effects: []itemEffect{
		{
			effectType: RecoverDiscardedO2,
			value:      3,
		},
	},
}

var Sonar = item{
	name: "Sonar",
	effects: []itemEffect{
		{
			effectType: ReorderNextO2Cards,
			value:      3,
		},
	},
}

var commonItemCards = []itemCard{
	{
		genericCard: genericCard{
			Name: "Common Card 1",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 2",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 3",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 4",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 5",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 6",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 7",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Common Card 8",
			Type: ItemType,
		},
		rarity: Common,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
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
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 2",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 3",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 4",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 5",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Uncommon Card 6",
			Type: ItemType,
		},
		rarity: Uncommon,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
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
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Rare Card 2",
			Type: ItemType,
		},
		rarity: Rare,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Rare Card 3",
			Type: ItemType,
		},
		rarity: Rare,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Rare Card 4",
			Type: ItemType,
		},
		rarity: Rare,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
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
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
	{
		genericCard: genericCard{
			Name: "Legendary Card 2",
			Type: ItemType,
		},
		rarity: Legendary,
		items: map[int]item{
			1: FlashLight,
			2: FlashLight,
			3: FlashLight,
			4: FlashLight,
			5: FlashLight,
			6: FlashLight,
			7: FlashLight,
			8: FlashLight,
			9: FlashLight,
		},
	},
}
