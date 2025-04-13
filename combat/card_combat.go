package combat

import (
	"my_test/log"
	"my_test/util"
)

const (
	CARD_COUNT  = 3
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
	deck      []*Card
	Hand      []*Card
	discard   []*Card
	remove    []*Card
	maxCard   int
	energy    int
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
		energy:      ENERGY_INIT,
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

func (c *CardCombat) StartTurn(action Action) {
	cardsToUse := []*Card{}
	for _, name := range action.Cards {
		card := c.GetCard(name)
		if card != nil {
			cardsToUse = append(cardsToUse, card)
		}
	}

	if len(cardsToUse) > 0 {
		c.Use(cardsToUse, c.Actors(), c.Enemies())
	}

	for _, discardName := range action.Discards {
		card := c.GetCard(discardName)
		if card != nil {
			c.DiscardCard(card)
		}
	}
}

func (c *CardCombat) EndTurn(action Action) {
	cardsToUse := []*Card{}
	for _, name := range action.Cards {
		card := c.GetCard(name)
		if card != nil {
			cardsToUse = append(cardsToUse, card)
		}
	}

	if len(cardsToUse) > 0 {
		c.Use(cardsToUse, c.Actors(), c.Enemies())
	}

	for _, discardName := range action.Discards {
		card := c.GetCard(discardName)
		if card != nil {
			c.DiscardCard(card)
		}
	}
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

func (c *CardCombat) CombatOnce(attacker Combatable, defender Combatable, isActorAttacker bool) CombatOnceResult {
	attacker.OnAttack(defender)
	damage := c.cacDamage(attacker, defender)
	if c.shouldDodge(attacker, defender) {
		log.Info("%s dodge damage %d on %s", attacker.GetName(), damage, defender.GetName())
		return CombatOnceResult{
			attackerDead: false,
			defenderDead: false,
		}
	}
	defender.OnDamage(damage, attacker)
	if isActorAttacker {
		c.actorCastDamage += damage
	} else {
		c.actorIncurDamage += damage

	}
	log.Info("%s cast damage %d on %s", attacker.GetName(), damage, defender.GetName())
	return CombatOnceResult{
		attackerDead: !attacker.IsAlive(),
		defenderDead: !defender.IsAlive(),
	}
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
