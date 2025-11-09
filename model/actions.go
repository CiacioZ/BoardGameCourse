package model

type actionType string

const (
	Breath    actionType = "BREATH"
	Explore   actionType = "EXPLORE"
	Dive      actionType = "DIVE"
	CalmDown  actionType = "CALM_DOWN"
	Ascend    actionType = "ASCEND"
	Distract  actionType = "DISTRACT"
	UseObject actionType = "USE_OBJECT"
)

type actionParam string

const (
	ExploreTime  actionParam = "EXPLORE_TIME"
	DiveLevels   actionParam = "DIVE_LEVELS"
	AscendLevels actionParam = "ASCEND_LEVELS"
	ItemToUse    actionParam = "ITEM_TO_USE"
)

type action struct {
	actionType actionType
	params     map[actionParam]int
}

func NewAction(actionType actionType, params map[actionParam]int) action {
	return action{
		actionType: actionType,
		params:     params,
	}
}
