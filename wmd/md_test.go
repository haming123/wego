package wmd

import (
	"os"
	"testing"
)

//go test -v -run=none -bench="BenchmarkMarshalHtml" -benchmem
func BenchmarkMarshalHtml(b *testing.B) {
	file_name := "./demo.md"
	input, err := os.ReadFile(file_name)
	if err != nil {
		b.Error(err)
		return
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MarshalHtml(input)
	}
	b.StopTimer()
}

