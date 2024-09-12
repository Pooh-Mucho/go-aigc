package aigc

import (
	"github.com/Pooh-Mucho/go-aigc/internal"
	"math"
)

var Tokenizer = tokenizer{}

type tokenizer struct{}

func (t *tokenizer) FastEstimate(s string) int {
	if len(s) == 0 {
		return 0
	}

	var (
		ac int // alphabet count
		nc int // number count
		sc int // symbol count
		uc int // unicode count
	)

	var buf = internal.UnsafeStringToBytes(s)
	for _, c := range buf {
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
			ac++
		} else if c >= '0' && c <= '9' {
			nc++
		} else if c <= 0x80 {
			sc++
		} else {
			uc++
		}
	}

	return int(math.Ceil(float64(ac)*0.192) + math.Ceil(float64(nc)*0.423) +
		math.Ceil(float64(sc)*0.5) + math.Ceil(float64(uc)*0.481))
}
