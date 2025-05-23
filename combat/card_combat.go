package combat

import (
	"errors"
	"my_test/event"
	"my_test/log"
	"my_test/push"
	"my_test/util"
	"slices"
	"time"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

const (
	UI_ACTOR_HP = iota
	UI_ACTOR_MAX_HP
	UI_ENEMY_HP
	UI_ENEMY_MAX_HP
)

const (
	TYPE_BUFF = iota
	TYPE_DEBUFF
)

const (
	BUFF_VULNERABLE = "vulnerable"
	BUFF_WEAK       = "weak"
	BUFF_STRENGTH   = "strength"
	BUFF_ARMOR      = "armor"
	BUFF_POISON     = "poison"
)

type CardEffect struct {
	Effect string `json:"effect"`
	Value  any    `json:"value"`
}

type CardCareer struct {
	InitCards []*Card
}

type CardEvent struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Effects     []CardEffect `json:"effects"`
}

type TurnInfo struct {
	drawedCards    []*Card
	usedCards      []*Card
	discardCards   []*Card
	exhaustedCards []*Card
	costEnergy     int
}

type CardBonus struct {
	Bonus []string
}

type CardCombatParams struct {
	Cards       []*Card
	Actors      []*CardActor
	Enemies     []*CardEnemy
	ResourceDir string
	CardCombatDelegate
}

type CardCombat struct {
	delegate        CardCombatDelegate
	combatables     []Combatable
	actors          []*CardActor
	originalEnemies []*CardEnemy
	enemies         []*CardEnemy
	ai              *EnemyAI
	deck            []*Card
	hand            []*Card
	discard         []*Card
	exhaust         []*Card
	initCardCount   int
	initEnergy      int
	turnCount       int
	uiDirty         bool
	uiTimer         time.Ticker
	finish          bool
	hurtCount       int
	attackCount     int
	turnInfo        TurnInfo
}

func NewCardCombat(p *CardCombatParams) *CardCombat {
	c := &CardCombat{
		actors:          p.Actors,
		originalEnemies: p.Enemies,
		enemies:         p.Enemies,
		delegate:        p.CardCombatDelegate,
		ai:              NewEnemyAI(p.Enemies),
		combatables:     make([]Combatable, len(p.Actors)+len(p.Enemies)),
		deck:            make([]*Card, len(p.Cards)),
		initCardCount:   p.Actors[0].InitCardCount,
		initEnergy:      p.Actors[0].InitEnergy,
		uiDirty:         false,
		uiTimer:         *time.NewTicker(time.Millisecond * 300),
		finish:          false,
		hurtCount:       0,
		attackCount:     0,
	}
	i := 0
	for _, actor := range p.Actors {
		c.combatables[i] = actor
		i++
	}
	for _, enemy := range p.Enemies {
		c.combatables[i] = enemy
		i++
	}
	for i, card := range p.Cards {
		c.deck[i] = card.Copy()
	}

	return c
}

func (c *CardCombat) Start() {
	c.turnCount = 0
	c.StartTurn()
	push.PushAction("战斗开始")
	c.requestUpdateUI()

	go c.updateUI()
}

func (c *CardCombat) Enemies() []Combatable {
	result := make([]Combatable, len(c.enemies))
	for i, enemy := range c.enemies {
		result[i] = enemy
	}
	return result
}

func (c *CardCombat) Actors() []Combatable {
	result := make([]Combatable, len(c.actors))
	for i, actor := range c.actors {
		result[i] = actor
	}
	return result
}

func (c *CardCombat) Combatables() []Combatable {
	return c.combatables
}

func (c *CardCombat) StartTurn() {
	c.turnCount++
	c.EnableBuff()
	c.delegate.OnActorTurnStart()
	c.actors[0].Energy = c.initEnergy
	drawCount := c.initCardCount - len(c.hand)
	mutable.Shuffle(c.deck)
	c.DrawCards(drawCount)
	for _, actor := range c.actors {
		actor.UpdateBuffs()
	}
	c.turnInfo = TurnInfo{}
	c.PrepareIntent()
}

func (c *CardCombat) EnableBuff() {
	c.delegate.EnableBuff()
}

func (c *CardCombat) PrepareIntent() {
	c.delegate.TriggerPrepareIntent()
}

func (c *CardCombat) UseCard(idx int32, choosenIdx []int32, target int32) error {
	card := c.hand[idx]
	if card == nil {
		return errors.New("card not exist")
	}
	if !c.delegate.CanUse(card) {
		return errors.New(card.Name + " cannot be used")
	}

	var choosen []*Card
	if choosenIdx != nil && len(choosenIdx) > 0 {
		choosen = make([]*Card, len(choosenIdx))
		for i, idx := range choosenIdx {
			switch card.GetValue("choose_from") {
			case CARD_CHOOSE_FROM_HAND:
				choosen[i] = c.hand[idx]
			case CARD_CHOOSE_FROM_DRAW:
				choosen[i] = c.deck[idx]
			case CARD_CHOOSE_FROM_DISCARD:
				choosen[i] = c.discard[idx]
			case CARD_CHOOSE_FROM_EXHAUST:
				choosen[i] = c.exhaust[idx]
			}
		}
	}

	// when choose_count is 0, choosen is all cards but used card in hand
	if card.HasValue("choose_count") && card.GetValue("choose_count") == 0 {
		choosen = lo.Filter(c.hand, func(c *Card, i int) bool {
			return i != int(idx)
		})
	}

	if c.actors[0].Energy >= card.Cost {
		c.actors[0].Energy -= card.Cost
		c.Use(card, choosen, c.enemies[target])
		c.delegate.OnUseCard(card)
	}
	c.turnInfo.usedCards = append(c.turnInfo.usedCards, card)

	c.requestUpdateUI()
	return nil
}

func (c *CardCombat) DiscardCards(cards []int32) int {
	for _, idx := range cards {
		c.discard = append(c.discard, c.hand[idx])
		c.hand[idx] = nil
	}

	newHand := []*Card{}
	for _, v := range c.hand {
		if v != nil {
			newHand = append(newHand, v)
		}
	}
	c.hand = newHand

	c.requestUpdateUI()
	push.PushAction("discard %d cards", len(cards))
	return len(c.discard)
}

func (c *CardCombat) EndTurn() {
	c.delegate.OnActorTurnEnd()
	discardCount := len(c.hand) - c.initCardCount
	if discardCount > 0 {
		c.discard = append(c.discard, c.hand[c.initCardCount:]...)
		c.hand = c.hand[:c.initCardCount]
	}

	c.EnemyTurn()

	if !c.finish {
		c.StartTurn()
		c.requestUpdateUI()
		push.PushAction("end turn")
	}
}

func (c *CardCombat) EnemyAttack(enemy *CardEnemy) {
	actor := c.actors[0]
	enemy.Attack = enemy.GetValue("attack")
	damage := c.CastDamage(enemy, actor)
	if damage != 0 {
		c.hurtCount++
		c.delegate.OnEnenyDamage(enemy, damage)
	}
}

func (c *CardCombat) EnemyMultiAttack(enemy *CardEnemy) {
	times := enemy.GetValue("attack_times")
	for i := 0; i < times; i++ {
		if !c.finish {
			c.EnemyAttack(enemy)
		}
	}
}

func (c *CardCombat) EnemyAddArmor(enemy *CardEnemy) {
	armor := enemy.GetValue("defense")
	enemy.AddBuff(Buff{
		Name:  BUFF_ARMOR,
		Type:  TYPE_BUFF,
		Value: armor,
	})
	push.PushAction("%s 施加了护盾", enemy.GetName())

}

func (c *CardCombat) EnemyTurn() {
	c.delegate.OnEnemyTurnStart()
	for _, v := range c.enemies {
		v.UpdateBuffs()
		c.checkDead(v)
	}

	if c.finish {
		return
	}

	for _, enemy := range c.enemies {
		if c.finish {
			break
		}
		action := c.ai.GetAction(enemy)
		bindings := make(map[string]any)
		bindings["enemy"] = enemy
		push.PushAction("%s action: %s", enemy.Name, action.Action)
		c.delegate.TriggerEnemyAction(action.Action, bindings)
	}

	c.delegate.OnEnemyTurnEnd()
	c.requestUpdateUI()
	c.ai.onEnemyTurnFinish()
}

func (c *CardCombat) checkDead(cbt Combatable) {
	if !cbt.IsAlive() {
		cbt.OnDead(nil)
		c.removeCombatable(cbt)
	}
}

// TODO: 当uiDirty为true时，才触发timer
func (c *CardCombat) requestUpdateUI() {
	c.uiDirty = true
}

// TODO: 线程安全问题
func (c *CardCombat) updateUI() {
	for {
		<-c.uiTimer.C
		if c.uiDirty {
			ev := &event.CardUpdateUIEvent{
				Actor: make([]event.ActorCardUI, len(c.actors)),
				Enemy: make([]event.EnemyCardUI, len(c.enemies)),
				Deck: event.DeckUI{
					DrawCount:    len(c.deck),
					DiscardCount: len(c.discard),
					ExhaustCount: len(c.exhaust),
					HandCards:    c.getHandString(),
				},
			}
			for i := range c.actors {
				copier.Copy(&ev.Actor[i], &c.actors[i])
			}
			for i := range c.enemies {
				copier.Copy(&ev.Enemy[i], &c.enemies[i])
				copier.Copy(&ev.Enemy[i].Intent, c.ai.GetAction(c.enemies[i]))
			}
			push.PushEvent(*ev)
			c.uiDirty = false
		}
	}
}

func (c *CardCombat) CastDamage(attacker Combatable, defender Combatable) int {
	damage := c.CacDamage(attacker, defender)
	armorBuff := defender.GetBase().GetArmorBuff()
	if armorBuff != nil {
		armorValue := armorBuff.Value
		armorBuff.Value -= damage
		if armorBuff.Value < 0 {
			defender.GetBase().RemoveBuff(BUFF_ARMOR)
		}
		damage = max(damage-armorValue, 0)
	}
	if damage > 0 {
		c.checkDead(defender)
		defender.GetBase().Life -= damage
		defender.OnDamage(damage, attacker)
	}
	push.PushAction("%s 攻击了 %s 造成 %d 点伤害", attacker.GetName(), defender.GetName(), damage)
	return damage
}

func (c *CardCombat) CacDamage(attacker Combatable, defender Combatable) int {
	buffStrength := attacker.GetBase().GetBuffValue(BUFF_STRENGTH)
	damage := attacker.GetAttack() + attacker.GetBase().Strength + buffStrength
	if attacker.GetBase().HasBuff(BUFF_WEAK) {
		damage = damage * attacker.GetBase().WeakFactor / 100
	}
	if defender.GetBase().HasBuff(BUFF_VULNERABLE) {
		damage = damage * defender.GetBase().VulnerableFactor / 100
	}
	return max(damage, 0)
}

func (c *CardCombat) removeCombatable(combatable Combatable) {
	switch combatable.GetCombatType() {
	case ACTOR:
		for i, actor := range c.actors {
			if actor == combatable {
				c.actors = slices.Delete(c.actors, i, i+1)
				break
			}
		}
	case ENEMY:
		for i, enemy := range c.enemies {
			if enemy == combatable {
				c.delegate.OnEnemyDead(enemy)
				c.enemies = slices.Delete(c.enemies, i, i+1)
				break
			}
		}
	default:
		log.Info("unknown combatable type %d", combatable.GetCombatType())
	}

	if len(c.actors) == 0 {
		c.onCombatFinish(false)
	} else if len(c.enemies) == 0 {
		c.onCombatFinish(true)
	}

	for i, v := range c.combatables {
		if v == combatable {
			c.combatables = slices.Delete(c.combatables, i, i+1)
			return
		}
	}
}

func (c *CardCombat) onCombatFinish(win bool) {
	c.finish = true

	if win {
		c.delegate.OnWin()
	} else {
		c.delegate.OnLose()
	}

	push.PushAction("combat finish")
}

func (c *CardCombat) getHandString() []string {
	strs := make([]string, 0, len(c.hand))
	for _, card := range c.hand {
		strs = append(strs, card.Id)
	}
	return strs
}

func (c *CardCombat) UpgradeCards(cards []*Card) {
	for _, card := range cards {
		c.upgradeCard(card)
	}
}

func (c *CardCombat) upgradeCard(card *Card) {
	card.Name = card.Name + "+"
	c.delegate.UpgradeCardInCombat(card)
}

func (c *CardCombat) AddCard(name string) {
	card := c.delegate.GetCard(name)
	c.delegate.OnAddCard(card)
	c.hand = append(c.hand, card)
	c.requestUpdateUI()
}

func (c *CardCombat) DrawCards(n int) []*Card {
	cards := make([]*Card, 0, n)
	for i := 0; i < n; i++ {
		card := c.DrawCard()
		if card == nil {
			break
		}
		cards = append(cards, card)
	}
	c.turnInfo.drawedCards = append(c.turnInfo.drawedCards, cards...)

	return cards
}

func (c *CardCombat) DrawCard() *Card {
	if len(c.deck) == 0 {
		if len(c.discard) == 0 {
			return nil
		}
		c.deck = append(c.deck, c.discard...)
		c.discard = c.discard[:0]
		c.ShuffleDeck()
	}
	top := len(c.deck) - 1
	card := c.FetchFromDeck(top)
	c.delegate.OnDrawCard(card)
	c.hand = append(c.hand, card)
	return card
}

func (c *CardCombat) FetchFromDeck(index int) *Card {
	if index < 0 || index >= len(c.deck) {
		return nil
	}
	card := c.deck[index]
	c.deck = slices.Delete(c.deck, index, index+1)
	return card
}

func (c *CardCombat) UseDeckTop() bool {
	if len(c.deck) == 0 {
		return false
	}
	top := len(c.deck) - 1
	card := c.FetchFromDeck(top)
	ramdomEnemy := lo.Sample(c.enemies)
	c.Use(card, nil, ramdomEnemy)
	return true
}

func (c *CardCombat) PutOnDeck(cards []*Card) {
	for _, card := range cards {
		c.deck = append(c.deck, card)
	}
	c.requestUpdateUI()
}

func (c *CardCombat) FetchFromDiscard(indexes []int) []*Card {
	cards := make([]*Card, 0, len(indexes))
	for _, idx := range indexes {
		cards = append(cards, c.discard[idx])
		c.discard[idx] = nil
	}
	c.discard = slices.DeleteFunc(c.discard, func(c *Card) bool {
		return c == nil
	})
	return cards
}

func (c *CardCombat) DiscardCard(card *Card) {
	c.delegate.OnDiscardCard(card)
	for i, v := range c.hand {
		if v == card {
			c.hand = slices.Delete(c.hand, i, i+1)
			c.discard = append(c.discard, card)
			return
		}
	}
	c.turnInfo.discardCards = append(c.turnInfo.discardCards, card)
}

func (c *CardCombat) PutCardDiscard(card *Card) {
	c.discard = append(c.discard, card)
	c.turnInfo.discardCards = append(c.turnInfo.discardCards, card)
}

func (c *CardCombat) ExhaustCard(card *Card) {
	c.delegate.OnExhaustCard(card)
	c.exhaust = append(c.exhaust, card)
	c.turnInfo.exhaustedCards = append(c.turnInfo.exhaustedCards, card)
}

func (c *CardCombat) ShuffleDeck() {
	c.delegate.OnShuffle()
	for i := len(c.deck) - 1; i > 0; i-- {
		j := util.GetRandomInt(i + 1)
		c.deck[i], c.deck[j] = c.deck[j], c.deck[i]
	}
}

func (c *CardCombat) GetHandCards() []*Card {
	return c.hand
}

func (c *CardCombat) getTargetsFromRange(card *Card, target *CardEnemy) []*CardEnemy {
	targets := make([]*CardEnemy, 0)
	if card.Range == CARD_RANGE_SINGLE {
		targets = append(targets, target)
	} else if card.Range == CARD_RANGE_ALL {
		targets = c.enemies
	} else if card.Range == CARD_RANGE_RANDOM {
		target := lo.Sample(c.enemies)
		targets = append(targets, target)
	}
	return targets
}

func (c *CardCombat) Attack(card *Card, target *CardEnemy) {
	targets := c.getTargetsFromRange(card, target)
	for _, target := range targets {
		c.actors[0].Attack = card.GetValue("attack")
		damage := c.CastDamage(c.actors[0], target)
		if damage != 0 {
			c.attackCount++
		}
	}
	c.requestUpdateUI()
}

func (c *CardCombat) MultiAttack(card *Card, target *CardEnemy) {
	times := card.GetValue("attack_times")
	for i := 0; i < times; i++ {
		if !c.finish {
			c.Attack(card, target)
		}
	}
}

func (c *CardCombat) AddArmor(card *Card) {
	armor := card.GetValue("defense")
	c.actors[0].AddBuff(Buff{
		Name:  BUFF_ARMOR,
		Type:  EFFECT_TYPE_BUFF,
		Value: armor,
	})

	bindings := make(map[string]any)
	bindings["armor"] = armor
	c.delegate.TriggerTiming(TIMING_ADD_ARMOR, bindings)
	c.requestUpdateUI()
}

func (c *CardCombat) AddVulnerable(card *Card, target *CardEnemy) {
	targets := c.getTargetsFromRange(card, target)
	for _, target := range targets {
		value := card.GetValue("vulnerable")
		target.AddBuff(Buff{
			Name: BUFF_VULNERABLE,
			Type: TYPE_DEBUFF,
			Turn: value,
		})

		c.delegate.TriggerTiming(TIMING_ADD_DEBUFF, nil)
	}
	c.requestUpdateUI()
}

func (c *CardCombat) AddWeak(card *Card, target *CardEnemy) {
	targets := c.getTargetsFromRange(card, target)
	for _, target := range targets {
		value := card.GetValue("weak")
		target.AddBuff(Buff{
			Name: BUFF_WEAK,
			Type: TYPE_DEBUFF,
			Turn: value,
		})

		c.delegate.TriggerTiming(TIMING_ADD_DEBUFF, nil)
	}
	c.requestUpdateUI()
}

func (c *CardCombat) AddEnemyEffect(enemy *CardEnemy, timing string, rule string) {
	effect := &Effect{
		CasterType: ENEMY,
		CasterID:   enemy.Name,
		Rule:       rule,
	}
	effect.Timing = TimingStr2Int(timing)
	c.delegate.AddEnemyEffect(effect)
}

func (c *CardCombat) AddEnemyBuff(enemy *CardEnemy, timing string, rule string) {
	effect := &Effect{
		Type:       EFFECT_TYPE_BUFF,
		Enabled:    false,
		CasterType: ENEMY,
		CasterID:   enemy.Name,
		Rule:       rule,
	}
	effect.Timing = TimingStr2Int(timing)
	c.delegate.AddEnemyEffect(effect)
}

func (c *CardCombat) AddBuff(target Combatable, name string, typ, value int, turn int) {
	target.GetBase().AddBuff(Buff{
		Name:  name,
		Type:  typ,
		Value: value,
		Turn:  turn,
	})
}

func (c *CardCombat) Use(card *Card, choosen []*Card, target *CardEnemy) {
	for _, effect := range card.Effects {
		bindings := make(map[string]any)
		bindings["target"] = target
		bindings["card"] = card
		if choosen != nil && len(choosen) > 0 {
			bindings["choosen_cards"] = choosen
		}
		c.delegate.TriggerEffect(effect, bindings)
	}
	c.removeHandCard(card)
	c.discard = append(c.discard, card)
	c.turnInfo.costEnergy += card.Cost
}

func (c *CardCombat) removeHandCard(card *Card) {
	for i, v := range c.hand {
		if v == card {
			c.hand = append(c.hand[:i], c.hand[i+1:]...)
			return
		}
	}
}
