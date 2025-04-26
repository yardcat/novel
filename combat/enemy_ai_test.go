package combat

import (
	"testing"
)

func TestEnemyAI_onEnemyTurnFinish(t *testing.T) {
	s := []int{1, 2, 3}
	p := &s[0]
	s = append(s, 4, 5, 6) // 可能触发扩容
	*p = 30
	print(s)
}
