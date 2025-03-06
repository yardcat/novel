package repo

import (
	"my_test/user"
)

type Cell struct {
	Item Item
	num  int
}

type Repo struct {
	Owner    *user.User
	Items    []*Cell
	Capacity int
}

func NewRepo(user *user.User, capacity int) *Repo {
	return &Repo{
		Owner:    user,
		Items:    []*Cell{},
		Capacity: capacity,
	}
}

func (r *Repo) Add(item Item) bool {
	if len(r.Items) < r.Capacity {
		r.Items = append(r.Items, &Cell{
			Item: item,
			num:  1,
		})
		return true
	}
	return false
}

func (r *Repo) Remove(item Item) bool {
	for i, cell := range r.Items {
		if cell.Item == item {
			r.Items = append(r.Items[:i], r.Items[i+1:]...)
			return true
		}
	}
	return false
}
