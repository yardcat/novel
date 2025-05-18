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
		ActionValue string `json:"value"`
		Description string `json:"description"`
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

type CardUpdateShopUI struct {
	Cards   []string `json:"cards"`
	Potions []string `json:"potions"`
	Relics  []string `json:"relics"`
}

type CardUpdatePotion struct {
	Potions []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"potions"`
}

type CardUpdateRelic struct {
	Relics []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"relics"`
}

type ActionUpdateEvent struct {
	Action string `json:"action"`
}

type CardCombatWin struct {
	Bonus struct {
		Cards             []string `json:"cards"`
		CardChooseCount   int      `json:"cardChooseCount"`
		Potions           []string `json:"potions"`
		PotionChooseCount int      `json:"potionChooseCount"`
		Relics            []string `json:"relics"`
		RelicChooseCount  int      `json:"relicChooseCount"`
	} `json:"bonus"`
	NextFloor []int `json:"next_floor"`
}

type CardCombatLose struct {
	Result string `json:"result"`
}

type CardEnterRoomDone struct {
	Type int `json:"type"`
}
