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
		Label  string
		Number int
	}
}
