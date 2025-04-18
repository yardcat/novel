package event

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
	Events     []string `json:"events"`
	Cards      []string `json:"handCards"`
	DeckCount  int      `json:"deckCount"`
	ActorHP    int      `json:"actorHP"`
	ActorMaxHP int      `json:"actorMaxHP"`
	EnemyHP    int      `json:"enemyHP"`
	EnemyMaxHP int      `json:"enemyMaxHP"`
	Energy     int      `json:"energy"`
}

type CardChooseStartEvent struct {
	Event string
}

type CardChooseStartEventReply struct {
	Results map[string]any
	Status  string
}

type CardSendCards struct {
	Cards []int
}

type CardSendCardsReply struct {
	Results map[string]any
}

type CardTurnEndEvent struct {
}

type CardTurnEndEventReply struct {
	DiscardCount int
	Damage       int
	NextAction   int
	ActionValue  int
	HandCards    string
	ActorHP      int
	ActorMaxHP   int
	EnemyHP      int
	EnemyMaxHP   int
}

type CardUpdateHandEvent struct {
	Cards []string `json:"cards"`
}

type CardUpdateUIEvent struct {
	Element string `json:"element"`
	Value   string `json:"value"`
}
