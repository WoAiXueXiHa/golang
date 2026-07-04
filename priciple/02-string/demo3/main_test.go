package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// 基础配置：拼接 1000 个短字符串
const (
	loopCount = 1000
	subStr    = "go"
)

// 1. + 操作符
func BenchmarkPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s string
		for j := 0; j < loopCount; j++ {
			s += subStr // 每次都会产生新字符串，旧字符串变垃圾，高频触发内存拷贝
		}
	}
}

// 2. fmt.Sprintf
func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s string
		for j := 0; j < loopCount; j++ {
			s = fmt.Sprintf("%s%s", s, subStr) // 内部涉及接口反射和动态分配，最慢
		}
	}
}

// 3. bytes.Buffer
func BenchmarkBytesBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		for j := 0; j < loopCount; j++ {
			buf.WriteString(subStr)
		}
		_ = buf.String() // 最后一次性转换为 string
	}
}

// 4. strings.Builder
func BenchmarkStringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		for j := 0; j < loopCount; j++ {
			builder.WriteString(subStr)
		}
		_ = builder.String() // 底层通过 unsafe 转换，零拷贝指针，性能极高
	}
}

// 5. append (切片转字符串)
func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf []byte
		for j := 0; j < loopCount; j++ {
			buf = append(buf, subStr...)
		}
		_ = string(buf) // 这一步依然会发生一次内存拷贝
	}
}

// 6. strings.Join
func BenchmarkStringsJoin(b *testing.B) {
	// 先准备好切片数据
	slice := make([]string, loopCount)
	for i := 0; i < loopCount; i++ {
		slice[i] = subStr
	}

	b.ResetTimer() // 重置时间，扣除准备切片的耗时
	for i := 0; i < b.N; i++ {
		_ = strings.Join(slice, "") // 内部提前计算总长度并预分配内存，适合已知切片拼接
	}
}
