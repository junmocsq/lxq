package tools

import (
	"math/rand"
	"testing"
)

func TestRandStr(t *testing.T) {
	t.Log(RandStr(100))
	t.Log(rand.Intn(2))
}

func TestPassword(t *testing.T) {
	t.Log(Password("junmo", "csq"))
}

func BenchmarkRandStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStr(10)
	}
}

func BenchmarkPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Password("junmojunmojunmojunmojunmojunmojunmojunmojunmojunmo", "csqcsqcsqcsq")
	}
}
