package career

import "my_test/util"

const (
	Doctor = iota
	Teacher
	Programmer
	CareerTypeCount
)

type Career struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Attr        map[string]util.Value `json:"attributes"`
}
