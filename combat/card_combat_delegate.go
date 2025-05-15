package combat

type CardCombatDelegate interface {
	GetCard(name string) *Card
	CanUse(card *Card) bool
	TriggerEffect(effect *Effect, binding map[string]any)
	TriggerTiming(timing int, binding map[string]any)
	OnWin()
	OnLose()
	OnUseCard(card *Card)
	OnDiscardCard(card *Card)
	OnShuffle()
	OnRemoveCard(card *Card)
	OnDrawCard(card *Card)
	OnAddCard(card *Card)
	// OnAttack(enemy *CardEnemy, damage int)
	// OnHurt(enemy *CardEnemy, damage int)
	// OnKillEnemy(enemy *CardEnemy)
	// OnDead(enemy *CardEnemy)
	// OnCurse(enemy *CardEnemy)
	// OnStrickBack(enemy *CardEnemy)
	// OnGetAmor(armor int)
	// OnGetStrength(armor int)
	OnEnemyDead(enemy *CardEnemy)
	OnEnemyTurnStart()
	OnEnemyTurnEnd()
	OnEnenyDamage(enemy *CardEnemy, damage int)
	OnActorTurnEnd()
	OnActorTurnStart()
}
