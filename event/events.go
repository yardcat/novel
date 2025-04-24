package event

type CardStatus struct {
	Name     string `json:"name"`
	Life     int    `json:"HP"`
	MaxLife  int    `json:"maxHP"`
	Energy   int    `json:"energy"`
	Strength int    `json:"strength"`
	Defense  int    `json:"defense"`
	Statuses []struct {
		Buff  string `json:"buff"`
		Value int    `json:"value"`
	} `json:"statuses"`
}

type CardUI struct {
	Name     string `json:"name"`
	Life     int    `json:"HP"`
	MaxLife  int    `json:"maxHP"`
	Energy   int    `json:"energy"`
	Strength int    `json:"strength"`
	Defense  int    `json:"defense"`
	Statuses []struct {
		Buff  string `json:"buff"`
		Value int    `json:"value"`
	} `json:"statuses"`
}

type DeckUI struct {
	DrawCount    int      `json:"drawCount"`
	DiscardCount int      `json:"discardCount"`
	HandCards    []string `json:"handCards"`
	NextAction   int      `json:"nextAction"`
	ActionValue  int      `json:"actionValue"`
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
	Events    []string `json:"events"`
	Cards     []string `json:"handCards"`
	DeckCount int      `json:"deckCount"`
}

type CardChooseStartEvent struct {
	Event string
}

type CardChooseStartEventReply struct {
	Results     map[string]any
	ActorStatus CardStatus `json:"actorStatus"`
	EnemyStatus CardStatus `json:"enemyStatus"`
}

type CardSendCards struct {
	Cards []int
}

type CardSendCardsReply struct {
	Status       string     `json:"status"`
	DrawCount    int        `json:"drawCount"`
	DiscardCount int        `json:"discardCount"`
	ActorStatus  CardStatus `json:"actorStatus"`
	EnemyStatus  CardStatus `json:"enemyStatus"`
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
	NextAction   int
	ActionValue  int
	HandCards    []string   `json:"handCards"`
	ActorStatus  CardStatus `json:"actorStatus"`
	EnemyStatus  CardStatus `json:"enemyStatus"`
}

type CardUpdateHandEvent struct {
	Cards []string `json:"cards"`
}

type CardUpdateUIEvent struct {
	Actor []CardUI `json:"actorUI"`
	Enemy []CardUI `json:"enemyUI"`
	Deck  DeckUI   `json:"deckUI"`
}
