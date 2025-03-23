package equip

type Armor struct {
	Attrs []Attr
}

func (a *Armor) GetAttrs() []Attr {
	return a.Attrs
}
