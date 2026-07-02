package main

import (
	"fmt"
	"go_tour/package1"
	"go_tour/utils"
)

func init() {
	fmt.Println("init func1 in main")
}

func init() {
	fmt.Println("init func2 in main")
}

var MainVal1 = utils.TraceLog("init Mval1", package1.V1 + 10)
var MainVal2 = utils.TraceLog("init Mval2", package1.V2 + 10)

func main() {
	fmt.Println("main func in main")
}