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
<<<<<<< HEAD
		Item  string
		Count int
	}
}

type CollectEventReply struct {
	EnergyCost int
	Items      []struct {
		Item  string
		Count int
=======
		Item  string `json:"item"`
		Count int    `json:"count"`
>>>>>>> 015ea37 (add reply task)
	}
}

type CollectEventReply struct {
	EnergyCost int
	Items      []struct {
		Item  string `json:"item"`
		Count int    `json:"count"`
	} `json:"items"`
}
