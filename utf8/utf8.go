package utf8

import (
	"unicode/utf8"
)

const RunesCount = utf8.MaxRune + 1

// Bounds of contiguous rune ranges. Up stands for uppercase and Lo for a lowercase.
const (
	ASCIIFirst, ASCIILast rune = 33, 126
	N0, N9                rune = 48, 57
	AUp, ZUp              rune = 65, 90
	ALo, ZLo              rune = 97, 122
)
