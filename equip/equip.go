package equip

const (
	ARMOR = iota
	BRACELET
	GLOVE
	HELMET
	NECKLACE
	RING
	SHOE
	TROUSER
	WEAPON
	EQUIP_COUNT
)

type Equip interface {
	GetAttrs() []Attr
}
