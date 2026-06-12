package report

import (
	"strconv"
	"strings"
)

// FormatBRL formata um valor no padrão monetário brasileiro: 214442.18 -> "214.442,18".
func FormatBRL(v float64) string {
	neg := v < 0
	if neg {
		v = -v
	}

	s := strconv.FormatFloat(v, 'f', 2, 64) // ex: "214442.18"
	parts := strings.SplitN(s, ".", 2)
	intPart, dec := parts[0], parts[1]

	var b strings.Builder
	n := len(intPart)
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteByte(intPart[i])
	}

	res := b.String() + "," + dec
	if neg {
		res = "-" + res
	}
	return res
}
