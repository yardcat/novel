package item

import (
	"my_test/user"
)

type Cell struct {
	Item Item
	num  int
}

type Bag struct {
	Owner    *user.Player
	Items    []*Cell
	Capacity int
}

func NewBag(user *user.Player, capacity int) *Bag {
	return &Bag{
		Owner:    user,
		Items:    []*Cell{},
		Capacity: capacity,
	}
}

func (b *Bag) Add(item Item) bool {
	if len(b.Items) < b.Capacity {
		b.Items = append(b.Items, &Cell{
			Item: item,
			num:  1,
		})
		return true
	}
	return false
}

func (b *Bag) Remove(item Item) bool {
	for i, cell := range b.Items {
		if cell.Item == item {
			b.Items = append(b.Items[:i], b.Items[i+1:]...)
			return true
		}
	}
	return false
}
