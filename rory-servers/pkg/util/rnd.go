package util

import (
	"math/rand"
	"time"
)

func init() {
	// 實際上使用 globalRand 來產生隨機數就行了
	// 自訂義的 var RndSeed rand.Source = rand.NewSource(time.Now().UnixNano()) 會有線程安全問題
	// NewSource returns a new pseudo-random Source seeded with the given value.
	// Unlike the default Source used by top-level functions, this source is not
	// safe for concurrent use by multiple goroutines.
	rand.Seed(time.Now().UnixNano())
}

// RndInt 生成範圍內的隨機整數, 不包含max
func RndInt(min int, max int) int {
	v := rand.Intn(max - min)
	return min + v
}

// RndFloat 生成0,1範圍內的隨機小數
func RndFloat() float32 {
	return rand.Float32()
}

// IsChance 隨機並計算是否命中機率
func IsChance(hit int, max int) bool {
	rnd := RndInt(0, max)
	if rnd < hit {
		return true
	}
	return false
}
