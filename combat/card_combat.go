package combat

import (
	"encoding/json"
	"my_test/event"
	"my_test/log"
	"my_test/push"
	"my_test/util"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/copier"
)

type CardCombat struct {
	combatables []Combatable
	actors      []*Actor
	enemies     []*Enemy
	ai          *EnemyAI
	client      CombatClient
	Record
	cardMap   map[string]*Card
	careerMap map[string]*CardCareer
	eventMap  map[string]*CardEvent
	deck      []*Card
	hand      []*Card
	discard   []*Card
	remove    []*Card
	maxCard   int
	turnNum   int
	uiDirty   bool
	uiTimer   time.Ticker
}

type EnemyTurnResult struct {
	damage      int
	action      string
	actionValue int
	actorDead   bool
	enemyDead   bool
}

func NewCardCombat(p *CombatParams) *CardCombat {
	c := &CardCombat{
		actors:      p.Actors,
		enemies:     p.Enemies,
		client:      p.Client,
		ai:          NewEnemyAI(p.Enemies),
		combatables: make([]Combatable, len(p.Actors)+len(p.Enemies)),
		cardMap:     make(map[string]*Card),
		careerMap:   make(map[string]*CardCareer),
		deck:        make([]*Card, 0),
		maxCard:     CARD_INIT,
		uiDirty:     false,
		uiTimer:     *time.NewTicker(time.Millisecond * 500),
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

	err := c.loadData(p.Path.GetPath("card"))
	if err != nil {
		panic(err)
	}
	return c
}

func (c *CardCombat) Start() EnemyAction {
	c.turnNum = 0
	c.PrepareCard()
	c.StartTurn()
	push.PushAction("战斗开始")
	go c.UpdateUI()
	return c.ai.EnemyAction(c.enemies[0])
}

func (c *CardCombat) ChooseDefender(attacker Combatable) Combatable {
	return c.enemies[0]
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
	c.actors[0].Energy = ENERGY_INIT
	drawCount := c.maxCard - len(c.hand)
	c.DrawCard(drawCount)
	for _, actor := range c.actors {
		actor.UpdateStatus()
	}
	c.ai.PrepareAction(c.enemies[0], c.actors)
}

func (c *CardCombat) UseCards(cards []int) *event.CardSendCardsReply {
	reply := &event.CardSendCardsReply{}
	cardsToUse := []*Card{}
	for _, idx := range cards {
		cardsToUse = append(cardsToUse, c.hand[idx])
	}
	results := make(map[string]any)
	for _, card := range cardsToUse {
		if c.actors[0].Energy >= card.Cost {
			c.actors[0].Energy -= card.Cost
			c.Use(card, results)
		}
	}

	c.requestUpdateUI()

	return reply
}

func (c *CardCombat) DiscardCards(ev *event.CardDiscardCards) *event.CardDiscardCardsReply {
	for _, idx := range ev.Cards {
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

	reply := &event.CardDiscardCardsReply{
		DiscardCount: len(c.discard),
	}
	return reply
}

func (c *CardCombat) EndTurn(ev *event.CardTurnEndEvent) *event.CardTurnEndEventReply {
	reply := &event.CardTurnEndEventReply{}
	discardCount := len(c.hand) - c.maxCard
	if discardCount > 0 {
		c.discard = append(c.discard, c.hand[c.maxCard:]...)
		c.hand = c.hand[:c.maxCard]
	}
	result := c.EnemyTurn()
	reply.Damage = result.damage
	c.actors[0].OnDamage(result.damage, c.enemies[0])

	c.StartTurn()
	action := c.ai.EnemyAction(c.enemies[0])
	copier.Copy(reply, action)

	c.requestUpdateUI()

	return reply
}

func (c *CardCombat) EnemyTurn() *EnemyTurnResult {
	result := &EnemyTurnResult{}

	for _, actor := range c.actors {
		actor.UpdateStatus()
	}

	action := c.ai.EnemyAction(c.enemies[0])
	if action.Action == ENEMY_BEHAVIOR_ATTACK {
		damage := c.cacDamage(c.enemies[0], c.actors[0])
		armorStatus := c.actors[0].GetBase().Statuses[STATUS_ARMOR]
		if armorStatus != nil && armorStatus.Value <= 0 {
			c.actors[0].RemoveStatus(STATUS_ARMOR)
			c.requestUpdateUI()
		}
		result.damage = damage
		result.actorDead = c.actors[0].GetLife() <= 0
		result.enemyDead = c.enemies[0].GetLife() <= 0
		push.PushAction("%s 攻击了 %s 造成 %d 点伤害", c.enemies[0].GetName(), c.actors[0].GetName(), damage)
	} else if action.Action == ENEMY_BEHAVIOR_DEFEND {
		c.enemies[0].AddStatus(Status{
			Type:  STATUS_ARMOR,
			Value: 5,
			Turn:  2,
		})
		push.PushAction("%s 施加了护盾", c.enemies[0].GetName())
	}

	c.requestUpdateUI()
	c.ai.onEnemyTurnFinish()

	return result
}

// TODO: 当uiDirty为true时，才触发timer
func (c *CardCombat) requestUpdateUI() {
	c.uiDirty = true
}

// TODO: 线程安全问题
func (c *CardCombat) UpdateUI() {
	for {
		<-c.uiTimer.C
		if c.uiDirty {
			ev := &event.CardUpdateUIEvent{
				Actor: make([]event.CardUI, len(c.actors)),
				Enemy: make([]event.CardUI, len(c.enemies)),
				Deck: event.DeckUI{
					DrawCount:    len(c.deck),
					DiscardCount: len(c.discard),
					HandCards:    c.getHandString(),
					// NextAction:   c.nextAction,
					// ActionValue:  c.actionValue,
				},
			}
			for i := range c.actors {
				copier.Copy(&ev.Actor[i], &c.actors[i])
			}
			for i := range c.enemies {
				copier.Copy(&ev.Enemy[i], &c.enemies[i])
			}
			push.PushEvent(*ev)
			c.uiDirty = false
		}
	}
}

func (c *CardCombat) cacDamage(attacker Combatable, defender Combatable) int {
	attack := attacker.GetAttack() + attacker.GetBase().Strength
	defense := defender.GetDefense() + defender.GetBase().Defense
	damage := attack - defense
	armorStatus := defender.GetBase().Statuses[STATUS_ARMOR]
	var armor int = 0
	if armorStatus != nil {
		armor = armorStatus.Value
		armorStatus.Value -= damage
	}
	damage -= armor
	return max(damage, 0)
}

func (c *CardCombat) removeCombatable(combatable Combatable) {
	switch combatable.GetCombatType() {
	case ACTOR:
		for i, actor := range c.actors {
			if actor == combatable {
				c.actors = append(c.actors[:i], c.actors[i+1:]...)
				break
			}
		}
	case ENEMY:
		for i, enemy := range c.enemies {
			if enemy == combatable {
				c.enemies = append(c.enemies[:i], c.enemies[i+1:]...)
				break
			}
		}
	default:
		log.Info("unknown combatable type %d", combatable.GetCombatType())
	}
	for i, v := range c.combatables {
		if v == combatable {
			c.combatables = append(c.combatables[:i], c.combatables[i+1:]...)
			return
		}
	}
}

func (c *CardCombat) onCombatFinish() {
	log.Info("combat finish, turns %d, actor cast %d damaage, actor incur %d damage",
		c.turns, c.actorCastDamage, c.actorIncurDamage)
	for _, actor := range c.actors {
		result := CombatResult{LifeCost: c.actorIncurDamage}
		actor.OnCombatDone(result)
	}
}

func (c *CardCombat) GetCard(name string) *Card {
	return c.cardMap[name]
}

func (c *CardCombat) GetCardTurnInfo() *CardTurnInfo {
	info := &CardTurnInfo{
		Cards:        make([]string, len(c.hand)),
		DrawCount:    len(c.deck),
		DiscardCount: len(c.discard),
		RemoveCount:  len(c.remove),
		Energy:       c.actors[0].Energy,
	}
	for i, card := range c.hand {
		info.Cards[i] = card.Name
	}
	return info
}

func (c *CardCombat) PrepareCard() {
	career := c.careerMap["kongfu"]
	c.deck = append(c.deck, career.InitCards...)
	c.ShuffleDeck()
}

func (c *CardCombat) GenerateChooseEvents() []string {
	return []string{"strength", "max_health", "draw_card"}
}

func (c *CardCombat) HandleChooseEvents(ev string) *event.CardChooseStartEventReply {
	reply := &event.CardChooseStartEventReply{
		Results: make(map[string]any),
	}
	for _, effect := range c.eventMap[ev].Effects {
		c.handCardEffect(&effect, reply.Results)
	}
	c.requestUpdateUI()
	return reply
}

func (c *CardCombat) getHandString() []string {
	strs := make([]string, 0, len(c.hand))
	for _, card := range c.hand {
		strs = append(strs, card.Name)
	}
	return strs
}

func (c *CardCombat) AddCard(card *Card) {
	c.hand = append(c.hand, card)
}

func (c *CardCombat) DrawCard(n int) []*Card {
	cards := make([]*Card, 0, n)

	for i := 0; i < n; i++ {
		if len(c.deck) == 0 {
			if len(c.discard) == 0 {
				break
			}
			c.deck = append(c.deck, c.discard...)
			c.discard = make([]*Card, 0)
			c.ShuffleDeck()
		}
		card := c.deck[0]
		c.deck = c.deck[1:]
		c.hand = append(c.hand, card)
		cards = append(cards, card)
	}

	return cards
}

func (c *CardCombat) DiscardCard(card *Card) {
	for i, v := range c.hand {
		if v == card {
			c.hand = append(c.hand[:i], c.hand[i+1:]...)
			c.discard = append(c.discard, card)
			return
		}
	}
}

func (c *CardCombat) RemoveCard(card *Card) {
	c.remove = append(c.remove, card)
}

func (c *CardCombat) ShuffleDeck() {
	for i := len(c.deck) - 1; i > 0; i-- {
		j := util.GetRandomInt(i + 1)
		c.deck[i], c.deck[j] = c.deck[j], c.deck[i]
	}
}

func (c *CardCombat) EffectFromString(effect string) int {
	switch effect {
	case "damage":
		return DAMAGE
	case "add_card":
		return ADD_CARD
	case "multi_damage":
		return MULTI_DAMAGE
	case "damage_defense":
		return DAMAGE_DEFENSE
	case "vulnerable":
		return VULNERABLE
	case "defend":
		return DEFEND
	case "weak":
		return WEAK
	case "strength":
		return STRENGTH
	case "heal":
		return HEAL
	case "max_health":
		return MAX_HEALTH
	default:
		return 0
	}
}

func (c *CardCombat) Use(card *Card, results map[string]any) {
	for _, effect := range card.Effects {
		c.handCardEffect(&effect, results)
	}
	c.removeHandCard(card)
	c.discard = append(c.discard, card)
}

func (c *CardCombat) handCardEffect(effect *CardEffect, results map[string]any) {
	switch c.EffectFromString(effect.Effect) {
	case DAMAGE:
		value := util.Anytoi(effect.Value)
		c.actors[0].Attack = value
		damage := c.cacDamage(c.actors[0], c.enemies[0])
		armorStatus := c.enemies[0].GetBase().Statuses[STATUS_ARMOR]
		if armorStatus != nil && armorStatus.Value <= 0 {
			c.enemies[0].RemoveStatus(STATUS_ARMOR)
			c.requestUpdateUI()
		}
		c.enemies[0].OnDamage(damage, c.actors[0])
		if damage > 0 {
			results["enemyHP"] = c.enemies[0].GetLife()
		}
	case VULNERABLE:
		for _, enemy := range c.enemies {
			enemy.AddStatus(Status{
				Type: STATUS_VULNERABLE,
				Turn: util.Anytoi(effect.Value),
			})
		}
	case DEFEND:
		for _, actor := range c.actors {
			actor.AddStatus(Status{
				Type:  STATUS_ARMOR,
				Value: util.Anytoi(effect.Value),
				Turn:  1,
			})
		}
	case MULTI_DAMAGE:
		n := util.Anytoi(effect.Value)
		for i := 0; i < n; i++ {
			for _, enemy := range c.enemies {
				enemy.OnDamage(util.Anytoi(effect.Value), nil)
			}
		}
	case DAMAGE_DEFENSE:
		for _, enemy := range c.enemies {
			enemy.OnDamage(util.Anytoi(effect.Value), nil)
		}
		for _, actor := range c.actors {
			actor.AddStatus(Status{
				Type: STATUS_ARMOR,
				Turn: util.Anytoi(effect.Value),
			})
		}
	case WEAK:
		for _, enemy := range c.enemies {
			enemy.AddStatus(Status{
				Type: STATUS_WEAK,
				Turn: util.Anytoi(effect.Value),
			})
		}
	case STRENGTH:
		value := util.Anytoi(effect.Value)
		c.actors[0].Strength += value
		results["strength"] = c.actors[0].Strength
	case HEAL:
		life := c.actors[0].Life + util.Anytoi(effect.Value)
		c.actors[0].Life = max(life, c.actors[0].MaxLife)
		results["actorHP"] = c.actors[0].Life
	case MAX_HEALTH:
		c.actors[0].MaxLife += util.Anytoi(effect.Value)
		results["actorMaxHP"] = c.actors[0].MaxLife
		c.actors[0].Life = c.actors[0].Life + util.Anytoi(effect.Value)
		results["actorHP"] = c.actors[0].Life
	case ADD_CARD:
		cards := effect.Value.([]any)
		for _, card := range cards {
			c.AddCard(c.GetCard(card.(string)))
		}
		push.PushEvent(event.CardUpdateHandEvent{Cards: c.getHandString()})
	}
}

func (c *CardCombat) removeHandCard(card *Card) {
	for i, v := range c.hand {
		if v == card {
			c.hand = append(c.hand[:i], c.hand[i+1:]...)
			return
		}
	}
}

func (c *CardCombat) loadData(dir string) error {
	err := c.loadCard(filepath.Join(dir, "card.json"))
	if err != nil {
		return err
	}
	err = c.loadCareer(filepath.Join(dir, "career.json"))
	if err != nil {
		return err
	}
	err = c.loadEvent(filepath.Join(dir, "event.json"))
	if err != nil {
		return err
	}
	return nil
}

func (c *CardCombat) loadCard(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &c.cardMap)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardCombat) loadCareer(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	adapter := map[string]struct {
		Cards []string `json:"init_cards"`
	}{}
	err = json.Unmarshal(data, &adapter)
	for k, v := range adapter {
		c.careerMap[k] = &CardCareer{
			InitCards: make([]*Card, 0, len(v.Cards)),
		}
		for _, card := range v.Cards {
			c.careerMap[k].InitCards = append(c.careerMap[k].InitCards, c.cardMap[card])
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *CardCombat) loadEvent(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &c.eventMap)
	if err != nil {
		return err
	}
	return nil
}
