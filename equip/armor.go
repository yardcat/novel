package equip

type Armor struct {
	Name  string
	Attrs []Attr
}

func (a *Armor) GetAttrs() []Attr {
	return a.Attrs
}
