package world

type Bag struct {
	Items   map[int]int
	capcity int
}

func NewBag() *Bag {
	return &Bag{
		Items:   make(map[int]int),
		capcity: 10,
	}
}

func (b *Bag) Add(item int) bool {
	if len(b.Items) >= b.capcity {
		return false
	}
	_, ok := b.Items[item]
	if !ok {
		b.capcity++
	}
	b.Items[item]++
	return true
}

func (b *Bag) GetCount(item int) int {
	return b.Items[item]
}

func (b *Bag) GetCapcity() int {
	return b.capcity
}

func (b *Bag) Remove(item int) bool {
	_, ok := b.Items[item]
	if !ok {
		return false
	}
	b.Items[item]--
	if b.Items[item] == 0 {
		delete(b.Items, item)
		b.capcity--
	}
	return true
}
