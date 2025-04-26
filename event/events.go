package event

type CardUI struct {
	Name     string `json:"name"`
	Life     int    `json:"HP"`
	MaxLife  int    `json:"maxHP"`
	Energy   int    `json:"energy"`
	Strength int    `json:"strength"`
	Defense  int    `json:"defense"`
	Statuses map[int]struct {
		Type  int `json:"type"`
		Value int `json:"value"`
		Turn  int `json:"turn"`
	} `json:"buffs"`
}

type DeckUI struct {
	DrawCount    int      `json:"drawCount"`
	DiscardCount int      `json:"discardCount"`
	HandCards    []string `json:"handCards"`
}

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

type CardStartEvent struct {
	Difficulty string
}

type CardStartEventReply struct {
	Events      []string `json:"events"`
	Cards       []string `json:"handCards"`
	DeckCount   int      `json:"deckCount"`
	Action      string   `json:"action"`
	ActionValue string   `json:"actionValue"`
}

type CardChooseStartEvent struct {
	Event string
}

type CardChooseStartEventReply struct {
	Results map[string]any
}

type CardSendCards struct {
	Cards []int
}

type CardSendCardsReply struct {
	Status       string `json:"status"`
	DrawCount    int    `json:"drawCount"`
	DiscardCount int    `json:"discardCount"`
}

type CardDiscardCards struct {
	Cards []int
}

type CardDiscardCardsReply struct {
	DiscardCount int `json:"discardCount"`
}

type CardTurnEndEvent struct {
}

type CardTurnEndEventReply struct {
	DrawCount    int `json:"drawCount"`
	DiscardCount int `json:"discardCount"`
	Damage       int
	Action       string   `json:"action"`
	ActionValue  int      `json:"actionValue"`
	HandCards    []string `json:"handCards"`
}

type CardUpdateHandEvent struct {
	Cards []string `json:"cards"`
}

type CardUpdateUIEvent struct {
	Actor []CardUI `json:"actorUI"`
	Enemy []CardUI `json:"enemyUI"`
	Deck  DeckUI   `json:"deckUI"`
}

type ActionUpdateEvent struct {
	Action string `json:"action"`
}
