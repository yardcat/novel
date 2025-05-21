package main

import (
	"encoding/json"
	"fmt"
	"my_test/combat"
	"my_test/log"
	"my_test/util"
	"os"
	"path/filepath"

	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
)

var (
	rootDir = ""
)

func loadRules() string {
	ruleData := ""
	ruleFile := filepath.Join(rootDir, "world", "island", "data", "card", "card_upgrade.rule")
	data, err := os.ReadFile(ruleFile)
	if err != nil {
		panic(err)
	}
	ruleData += string(data)
	return ruleData
}

func main() {
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}
	if rootDir == "" {
		panic("rootDir is empty, check your args")
	}

	jsonPath := filepath.Join(rootDir, "world", "island", "data", "card", "card.json")
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}
	cardMap := make(map[string]*combat.Card)
	err = json.Unmarshal(jsonData, &cardMap)
	if err != nil {
		panic(err)
	}

	dataContext := context.NewDataContext()
	dataContext.Add("log", fmt.Println)
	dataContext.Add("Sprintf", fmt.Sprintf)
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	err = ruleBuilder.BuildRuleFromString(loadRules())
	if err != nil {
		panic("load rule error")
	}
	engine := engine.NewGengine()

	newCardMap := make(map[string]*combat.Card)
	for id, card := range cardMap {
		upgradeRule := id + "_upgrade"
		dataContext.Add("card", card)
		exist := ruleBuilder.IsExist([]string{upgradeRule})
		if exist == nil || len(exist) == 0 || !exist[0] {
			log.Error("%s is not exist", upgradeRule)
			continue
		} else {
			log.Info("%s is upgraded", id)
		}
		err = engine.ExecuteSelectedRules(ruleBuilder, []string{upgradeRule})
		if err != nil {
			panic(err)
		}
		card.Description = util.FormatString(card.Description, card.Values)
		card.Name += "+"
		newId := id + "+"
		newCardMap[newId] = card
	}

	outData, err := json.Marshal(newCardMap)
	if err != nil {
		panic(err)
	}
	outPath := filepath.Join(rootDir, "world", "island", "data", "card", "card_upgrade.json")
	err = os.WriteFile(outPath, outData, 0644)
	if err != nil {
		panic(err)
	}
	log.Info("card upgrade success, output to %s", outPath)
}
