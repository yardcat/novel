package event

type ActorCardUI struct {
	Name     string `json:"name"`
	Life     int    `json:"HP"`
	MaxLife  int    `json:"maxHP"`
	Energy   int    `json:"energy"`
	Strength int    `json:"strength"`
	Defense  int    `json:"defense"`
	Statuses []struct {
		Type  int `json:"type"`
		Value int `json:"value"`
		Turn  int `json:"turn"`
	} `json:"buffs"`
}

type EnemyCardUI struct {
	Name     string `json:"name"`
	Life     int    `json:"HP"`
	MaxLife  int    `json:"maxHP"`
	Strength int    `json:"strength"`
	Defense  int    `json:"defense"`
	Intent   struct {
		Action      string `json:"action"`
		ActionValue int    `json:"value"`
		Description int    `json:"description"`
		Target      int    `json:"target"`
	} `json:"intent"`
	Statuses []struct {
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

type CardUpdateUIEvent struct {
	Actor []ActorCardUI `json:"actorUI"`
	Enemy []EnemyCardUI `json:"enemyUI"`
	Deck  DeckUI        `json:"deckUI"`
}

type ActionUpdateEvent struct {
	Action string `json:"action"`
}

type CardCombatWin struct {
	Bonus []string `json:"bonus"`
}

type CardCombatLose struct {
	Result string `json:"result"`
}
