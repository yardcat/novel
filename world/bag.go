package world

import (
	"fmt"
	"strings"
)

type Bag struct {
	Items   map[Item]int
	capcity int
}

func NewBag() *Bag {
	return &Bag{
		Items:   make(map[Item]int),
		capcity: 10,
	}
}

func (b *Bag) Add(item Item, count int) bool {
	if len(b.Items) >= b.capcity {
		return false
	}
	_, ok := b.Items[item]
	if !ok {
		b.capcity++
	}
	b.Items[item] += count
	return true
}

func (b *Bag) GetCount(item Item) int {
	return b.Items[item]
}

func (b *Bag) GetCapcity() int {
	return b.capcity
}

func (b *Bag) Remove(item Item) bool {
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

func (b *Bag) ToJson() string {
	res := `{"items": [`
	line_arr := []string{}
	for item := range b.Items {
		line := fmt.Sprintf(`{"name": "%s", "count": %d}`, item.GetName(), b.Items[item])
		line_arr = append(line_arr, line)
	}
	res += strings.Join(line_arr, ",")
	res += `]}`
	return res
}
