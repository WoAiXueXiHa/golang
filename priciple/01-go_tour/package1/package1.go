package package1

import (
	"fmt"
	"go_tour/package2"
	"go_tour/utils"
)

var V1 = utils.TraceLog("init package1 val1", package2.Val1 + 10)
var V2 = utils.TraceLog("init package1 val2", package2.Val2 + 10)


func init() {
	fmt.Println("init func in package1")
}