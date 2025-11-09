package model

type ActionType string

const (
	Breath    ActionType = "BREATH"
	Explore   ActionType = "EXPLORE"
	Dive      ActionType = "DIVE"
	CalmDown  ActionType = "CALM_DOWN"
	Ascend    ActionType = "ASCEND"
	Distract  ActionType = "DISTRACT"
	UseObject ActionType = "USE_OBJECT"
)

type ActionParam string

const (
	EploreTime   ActionParam = "EXPLORE_TIME"
	DiveLevels   ActionParam = "DIVE_LEVELS"
	AscendLevels ActionParam = "ASCEND_LEVELS"
	ObjectToUse  ActionParam = "OBJECT_TO_USE"
)

type action struct {
	actionType ActionType
	params     map[ActionParam]int
}

func NewAction(actionType ActionType, params map[ActionParam]int) action {
	return action{
		actionType: actionType,
		params:     params,
	}
}
