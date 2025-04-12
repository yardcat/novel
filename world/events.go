package world

type ChangeStatusEvent struct {
	Type  string
	Value int
}

type ChangeEnvEvent struct {
	Type  string
	Value string
}

type BonusEvent struct {
	Item  string
	Count int
}

type CollectEvent struct {
	Items []struct {
		Item  string `json:"item"`
		Count int    `json:"count"`
	}
}

type CollectEventReply struct {
	EnergyCost int
	Items      []struct {
		Item  string `json:"item"`
		Count int    `json:"count"`
	} `json:"items"`
}

type CombatWinEvent struct {
	Result string
}

type StartCardEvent struct {
	Event string
}

type StartCardEventReply struct {
	Card  string
	Event []string
}

type ChooseStartEvent struct {
	Card string
}

type ChooseStartEventReply struct {
	Card string
}

type CardTurnStartEvent struct {
	Card string
}

type CardTurnStartEventReply struct {
	Card string
}

type CardTurnEndEvent struct {
	Card string
}

type CardTurnEndEventReply struct {
	Card string
}
