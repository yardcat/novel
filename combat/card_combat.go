package combat

import (
	"my_test/event"
	"my_test/log"
	"my_test/push"
	"my_test/util"
	"strings"
)

const (
	CARD_COUNT  = 5
	ENERGY_INIT = 3
)

type CardCombat struct {
	combatables []Combatable
	actors      []*Actor
	enemies     []*Enemy
	client      CombatClient
	Record
	cardMap   map[string]*Card
	careerMap map[string]*CardCareer
	eventMap  map[string]*CardEvent
	deck      []*Card
	Hand      []*Card
	discard   []*Card
	remove    []*Card
	maxCard   int
	Energy    int
	turnNum   int
}

type EnemyTurnResult struct {
	damage     int
	nextAction int
	actorDead  bool
	enemyDead  bool
}

type Action struct {
	Cards    []string
	Discards []string
}

func NewCardCombat(p *CombatParams) *CardCombat {
	c := &CardCombat{
		actors:      p.Actors,
		enemies:     p.Enemies,
		client:      p.Client,
		combatables: make([]Combatable, len(p.Actors)+len(p.Enemies)),
		cardMap:     make(map[string]*Card),
		careerMap:   make(map[string]*CardCareer),
		deck:        make([]*Card, 0),
		maxCard:     CARD_COUNT,
		Energy:      ENERGY_INIT,
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

func (c *CardCombat) Start(difficuty string) error {
	log.Info("start card, difficulty:%s", difficuty)
	c.PrepareCard()
	c.turnNum = 0
	return nil
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
	c.Energy = ENERGY_INIT
	drawCount := c.maxCard - len(c.Hand)
	c.DrawCard(drawCount)
}

func (c *CardCombat) UseCards(cards []int) *event.CardSendCardsReply {
	reply := &event.CardSendCardsReply{}
	cardsToUse := []*Card{}
	for _, idx := range cards {
		cardsToUse = append(cardsToUse, c.Hand[idx])
	}
	results := make(map[string]any)
	for _, card := range cardsToUse {
		c.Use(card, results)
	}
	reply.Results = results
	return reply
}

func (c *CardCombat) EndTurn(ev *event.CardTurnEndEvent) *event.CardTurnEndEventReply {
	reply := &event.CardTurnEndEventReply{}
	discardCount := len(c.Hand) - c.maxCard
	if discardCount > 0 {
		reply.DiscardCount = discardCount
		c.Hand = c.Hand[:c.maxCard]
	}
	result := c.EnemyTurn()
	reply.Damage = result.damage
	reply.NextAction = result.nextAction

	c.StartTurn()
	reply.HandCards = strings.Join(c.getHandString(), ",")
	reply.ActorHP = c.actors[0].Life
	reply.ActorMaxHP = c.actors[0].MaxLife
	reply.EnemyHP = c.enemies[0].Life
	reply.EnemyMaxHP = c.enemies[0].MaxLife
	return reply
}

func (c *CardCombat) EnemyTurn() *EnemyTurnResult {
	result := &EnemyTurnResult{}
	result.damage = c.cacDamage(c.enemies[0], c.actors[0])
	result.nextAction = 0
	result.actorDead = c.actors[0].GetLife() <= 0
	result.enemyDead = c.enemies[0].GetLife() <= 0
	return result
}

func (c *CardCombat) UpdateUI(element string, value any) {
	push.PushEvent(event.CardUpdateUIEvent{})
}

func (c *CardCombat) ChooseAttacker() Combatable {
	fast := MAX_STEP
	fast_idx := 0
	for i, comb := range c.combatables {
		speed := (MAX_STEP - comb.GetBase().AttackStep) / float64(comb.GetAttackSpeed())
		if speed < fast {
			fast = speed
			fast_idx = i
		}
	}
	for _, comb := range c.combatables {
		comb.GetBase().AttackStep += float64(comb.GetAttackSpeed()) * fast
	}
	c.combatables[fast_idx].GetBase().AttackStep = 0
	return c.combatables[fast_idx]
}

func (c *CardCombat) cacDamage(attacker Combatable, defender Combatable) int {
	attack := attacker.GetAttack()
	defense := defender.GetDefense()
	damage_reduce_factor := 0.0
	damage := int(float64(attack-defense) * (1 - damage_reduce_factor))
	return max(damage, 0)
}

func (c *CardCombat) shouldDodge(_ Combatable, defender Combatable) bool {
	randomNumber := util.GetRandomInt(100)
	dodge := defender.GetDodge()
	return randomNumber < dodge
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
	log.Info("combat finish, turns %d, actor cast %d damaage, actor incur %d damage", c.turns, c.actorCastDamage, c.actorIncurDamage)
	for _, actor := range c.actors {
		result := CombatResult{LifeCost: c.actorIncurDamage}
		actor.OnCombatDone(result)
	}
}
