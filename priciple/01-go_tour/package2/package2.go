package package2

import (
	"fmt"
	"go_tour/utils"
)

var Val1 = utils.TraceLog("init package2 Val1", 20)
var Val2 = utils.TraceLog("init package2 Val2", 100)

func init() {
	fmt.Println("init func1 in package2")
}

func init() {
	fmt.Println("init func2 in package2")
}
