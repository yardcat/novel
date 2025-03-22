package career

const (
	Doctor = iota
	Teacher
	Programmer
	CareerTypeCount
)

type Career struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Attr        map[string]string `json:"attributes"`
}
