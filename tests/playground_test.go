//nolint:lll
package tests

import (
	"math"
	"testing"
)

func TestPlayground(t *testing.T) {
	t.SkipNow()
}

func BenchmarkPayload_SetABWithCopy(b *testing.B) {
	var payload Payload

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		payload = payload.SetWithCopy("some string value", math.MaxInt64, map[string]string{"some string value": "some other string value"}, []string{"asdf", "fdsa", "qwer", "bmdfjdal"})
	}
}

// BenchmarkPayload_SetABWithCopy-16    	1000000000	         0.2322 ns/op
// BenchmarkPayload_SetABWithPointer-16    	1000000000	         0.2347 ns/op

func BenchmarkPayload_SetABWithPointer(b *testing.B) {
	var payload Payload

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		payload.SetWithPointer("some string value", math.MaxInt64, map[string]string{"some string value": "some other string value"}, []string{"asdf", "fdsa", "qwer", "bmdfjdal"})
	}
}

type Payload struct {
	A string
	B int
	C map[string]string
	D []string
}

func (p Payload) SetWithCopy(a string, b int, c map[string]string, d []string) Payload {
	p.A = a
	p.B = b
	p.C = c
	p.D = d

	return p
}

func (p *Payload) SetWithPointer(a string, b int, c map[string]string, d []string) {
	p.A = a
	p.B = b
	p.C = c
	p.D = d
}
