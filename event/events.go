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
	DrawCount   int        `json:"drawCount"`
	ActorStatus CardStatus `json:"actorStatus"`
	EnemyStatus CardStatus `json:"enemyStatus"`
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
	DiscardCount int
	Damage       int
	NextAction   int
	ActionValue  int
	HandCards    string
	ActorStatus  CardStatus `json:"actorStatus"`
	EnemyStatus  CardStatus `json:"enemyStatus"`
}

type CardUpdateHandEvent struct {
	Cards []string `json:"cards"`
}

type CardUpdateUIEvent struct {
	Element string `json:"element"`
	Value   string `json:"value"`
}
