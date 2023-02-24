package cokacoin

import (
	"bytes"
	"testing"
)

func TestGenerateMask(t *testing.T) {
	data := []struct {
		In  int
		Out []byte
	}{
		{
			In:  0,
			Out: []byte{},
		},
		{
			In:  1,
			Out: []byte{0xf},
		},
		{
			In:  2,
			Out: []byte{0},
		},
	}

	for i := range data {
		out := GenerateMask(data[i].In)
		if !bytes.Equal(out, data[i].Out) {
			t.Errorf("Failed for %d, it should be %x but is %x\n", data[i].In, data[i].Out, out)
		}

	}
}

func BenchmarkDifficultHash(b *testing.B) {
	mask := GenerateMask(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DifficultHash(mask, "a", "b", i)

	}
}
