package stringutil

import (
	"fmt"
	"strings"
)

const maxInt = int(^uint(0) >> 1)

func Format(text string, args []any) string {
	if !strings.Contains(text, "{") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text) + len(text)/4)

	for i := 0; i < len(text); {
		if text[i] != '{' {
			j := i + 1
			for j < len(text) && text[j] != '{' {
				j++
			}
			b.WriteString(text[i:j])
			i = j
			continue
		}

		j := i + 1
		if j >= len(text) || text[j] < '0' || text[j] > '9' {
			b.WriteByte('{')
			i++
			continue
		}

		k := j
		for k < len(text) && text[k] >= '0' && text[k] <= '9' {
			k++
		}
		if k >= len(text) || text[k] != '}' {
			b.WriteByte('{')
			i++
			continue
		}

		n := 0
		for p := j; p < k; p++ {
			d := int(text[p] - '0')
			if n > (maxInt-d)/10 {
				panic("Invalid formatting placeholder")
			}
			n = n*10 + d
		}
		if n >= len(args) {
			panic("Invalid formatting placeholder")
		}

		b.WriteString(fmt.Sprint(args[n]))
		i = k + 1
	}

	return b.String()
}
