package main

import (
	"fmt"
	"math"
)

func traversal(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	// case1: 遍历顺序随机，同一 map 两次遍历顺序几乎必不同
	m := make(map[string]int, 8)
	for i := 0; i < 8; i++ {
		m[fmt.Sprintf("key%d", i)] = i
	}
	fmt.Println("case1 first :", traversal(m))
	fmt.Println("case1 second:", traversal(m))

	// case2: 读不存在的 key 返回零值，v, ok 可以区分
	v := m["nope"]
	_, ok := m["nope"]
	fmt.Println("case2 zero:", v, "ok:", ok)

	// case3: 删除已存在的 key 后 ok 变 false；删除不存在的 key 不报错
	delete(m, "key0")
	_, ok = m["key0"]
	fmt.Println("case3 after delete:", ok, "len:", len(m))
	delete(m, "never-existed")

	// case4: 遍历中删除是安全的，已删除的 key 不会再产生值
	for k := range m {
		delete(m, k)
	}
	fmt.Println("case4 len after iterate-delete:", len(m))

	// case5: NaN 键，NaN != NaN，所以每次写入都认为是"新键"
	nm := map[float64]int{}
	for i := 0; i < 3; i++ {
		nm[math.NaN()] = i
	}
	fmt.Println("case5 len:", len(nm))
	fmt.Println("case5 read NaN:", nm[math.NaN()])

	// case6: 0.0 与 -0.0 相等（且哈希相同），更新不新增
	zm := map[float64]int{}
	zm[0.0] = 1
	zm[math.Copysign(0, -1)] = 2
	fmt.Println("case6 len:", len(zm), "value:", zm[0.0])
}
