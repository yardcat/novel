package scene

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"my_test/user"
	"os"
	"strconv"
	"strings"
)

type SceneProto struct {
	id       int
	distance int
	monsters map[int]int
}

type MonsterProto struct {
	id      int
	name    string
	level   int
	life    int
	attack  int
	defense int
	dodge   int
}

var (
	scene_table   map[int]*SceneProto
	monster_table map[int]*MonsterProto
)

func CreateScene(id int, players []*user.Player) *LineScene {
	monsters := make([]*user.Enemy, 0)
	scene_proto := scene_table[id]
	for monster_id, num := range scene_proto.monsters {
		for j := 0; j < num; j++ {
			monster_proto := monster_table[monster_id]
			monsters = append(monsters, CreateMonster(monster_proto))
		}
	}
	return NewLineScene(players, monsters, scene_proto.distance)
}

func CreateMonster(proto *MonsterProto) *user.Monster {
	return user.NewMonster(proto.name, user.Property{
		Life:    proto.life,
		Attack:  proto.attack,
		Defense: proto.defense,
		Dodge:   proto.dodge,
	})
}

func LoadDataFromCSV() {
	loadMonsterTable()
	loadSceneTable()
}

func loadMonsterTable() {
	path := "data/monster.csv"
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("load monster table failed: ", err)
		return
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("load monster table failed: ", err)
		return
	}

	records_size := len(records)
	monster_table = make(map[int]*MonsterProto, records_size-1)
	for i := 1; i < records_size; i++ {
		record := records[i]
		proto := &MonsterProto{}
		proto.id, _ = strconv.Atoi(record[0])
		proto.name = record[1]
		proto.level, _ = strconv.Atoi(record[2])
		proto.life, _ = strconv.Atoi(record[3])
		proto.attack, _ = strconv.Atoi(record[4])
		proto.defense, _ = strconv.Atoi(record[5])
		proto.dodge, _ = strconv.Atoi(record[6])
		monster_table[proto.id] = proto
	}
	fmt.Println("load monster table success")
}

func loadSceneTable() {
	path := "data/scene.csv"
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("load scene table failed: ", err)
		return
	}
	defer f.Close()

	scene_table = make(map[int]*SceneProto)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		records := strings.Split(line, ",")
		proto := &SceneProto{monsters: make(map[int]int)}
		proto.id, _ = strconv.Atoi(records[0])
		proto.distance, _ = strconv.Atoi(records[1])
		for j := 2; j < len(records); j++ {
			pair := strings.Split(records[j], ":")
			monster_id, _ := strconv.Atoi(pair[0])
			num, _ := strconv.Atoi(pair[1])
			proto.monsters[monster_id] = num
		}
		scene_table[proto.id] = proto
	}
	fmt.Println("load scene table success")
}
