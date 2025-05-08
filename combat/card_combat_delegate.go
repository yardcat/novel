package combat

type CardCombatDelegate interface {
	GetCard(name string) *Card
	OnWin()
	OnLose()
	OnPlayCard(card *Card)
	OnDiscardCard(card *Card)
	OnShuffle()
	OnRemoveCard(card *Card)
	OnDrawCard(card *Card)
	OnAddCard(card *Card)
	OnEnemyDead(enemy *CardEnemy)
	OnEnemyTurnStart()
	OnEnemyTurnEnd()
	OnEnenyDamage(enemy *CardEnemy, damage int)
	OnActorTurnEnd()
	OnActorTurnStart()
}
