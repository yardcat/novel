package world

type MineResource struct {
	Store int `json:"store"`
	Rate  int `json:"rate"`
}

type MineEnemy struct {
	Rate int `json:"rate"`
}

type Mine struct {
	Resources map[string]MineResource `json:"resources"`
	Enemy     map[string]MineResource `json:"enemies"`
}

func NewMine() *Mine {
	return &Mine{}
}

func (m *Mine) Explore() {

}
