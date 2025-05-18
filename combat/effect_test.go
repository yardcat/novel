package combat

import (
	"fmt"
	"my_test/log"
	"testing"
	"time"

	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
)

// 定义想要注入的结构体
type User struct {
	Name string
	Age  int64
	Male bool
}

func (u *User) GetNum(i int64) int64 {
	return i
}

func (u *User) Print(s string) {
	fmt.Println(s)
}

func (u *User) Say() {
	fmt.Println("hello world")
}

// 定义规则
const rule1 = `
rule "name test" "i can"  salience 0
begin
		if 7 == User.GetNum(7){
			User.Age = User.GetNum(89767) + 10000000
			User.Print("6666")
		}else{
			User.Name = "yyyy"
		}
		return Sprintf("%s xxxxx %d",User.Name,User.Age)
end
`

func Test_Multi(t *testing.T) {
	user := &User{
		Name: "Calo",
		Age:  0,
		Male: true,
	}

	dataContext := context.NewDataContext()
	//注入初始化的结构体
	dataContext.Add("User", user)
	dataContext.Add("Sprintf", fmt.Sprintf)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	start1 := time.Now().UnixNano()
	//构建规则
	err := ruleBuilder.BuildRuleFromString(rule1) //string(bs)
	end1 := time.Now().UnixNano()

	log.Info("rules num:%d, load rules cost time:%d", len(ruleBuilder.Kc.RuleEntities), end1-start1)

	if err != nil {
		log.Error("err:%s ", err)
	} else {
		eng := engine.NewGengine()

		start := time.Now().UnixNano()
		//执行规则
		err := eng.Execute(ruleBuilder, true)
		println(user.Age)
		end := time.Now().UnixNano()
		if err != nil {
			log.Error("execute rule error: %v", err)
		}
		log.Info("execute rule cost %d ns", end-start)
		log.Info("user.Age=%d,Name=%s,Male=%t", user.Age, user.Name, user.Male)
	}
}
