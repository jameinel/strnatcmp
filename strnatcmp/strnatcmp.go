package strnatcmp

// isSpace is the Latin subset of unicode.IsSpace
func isSpace(c uint8) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

// nextDigit returns the value of the next character, and whether it is a digit
func nextDigit(s string, i int) (uint8, bool) {
	if i < len(s) {
		c := s[i]
		// we only handle latin digits, but this gets tricky if you
		// start considering non-latin numbers.
		return c, '0' <= c && c <= '9'
	}
	return 0, false
}

// compareDigits considers the current text as a number, and compares them.
// The longest run of digits wins, as it is clearly larger. However, if we get the same number of characters, then we have to consider
func compareDigits(a string, ai int, b string, bi int) int {
	var bias int = 0
	for {
		// We don't need to check if ai or bi is off the end of the
		// string, because one of them won't be a digit, which
		// guarantees a return from this function.
		char_a, a_digit := nextDigit(a, ai)
		char_b, b_digit := nextDigit(b, bi)
		if !a_digit && !b_digit {
			// both are no longer digits at the same time, so we
			// just return whatever bias we have observed.
			return bias
		} else if !a_digit {
			// no more digits in a, b has at least 1 extra digit
			return -1
		} else if !b_digit {
			// no more digits in b, a has at least 1 extra digit
			return 1
		} else if bias == 0 {
			// no bias yet, do the comparison
			if char_a < char_b {
				bias = -1
			} else if char_a > char_b {
				bias = 1
			}
		}
		ai++
		bi++
	}
}

// compareFractional handles when numbers start with a 0, to handle fractional
// decimals.  (1.01 is longer than 1.2, but comes before 1.2).
func compareFractional(a string, ai int, b string, bi int) int {
	// We just treat the first value to be larger as winning
	for {
		// We don't need to check if ai or bi is off the end of the
		// string, because one of them won't be a digit, which
		// guarantees a return from this function.
		char_a, a_digit := nextDigit(a, ai)
		char_b, b_digit := nextDigit(b, bi)
		if !a_digit && !b_digit {
			// Same length, and no differences
			return 0
		} else if !a_digit {
			// no more digits in a, b has at least 1 extra digit
			return -1
		} else if !b_digit {
			// no more digits in b, a has at least 1 extra digit
			return 1
		} else if char_a < char_b {
			return -1
		} else if char_a > char_b {
			return 1
		}
		ai++
		bi++
	}
}

// Strnatcmp compares two strings in "natural" order.
// It ignores whitespace and where it encounters digits, it sorts them
// numerically, rather than alphabetically. So
// foo10 comes after foo1 and foo2, rather than between them.
func Compare(a, b string) int {
	ai := 0
	bi := 0
	for {
		char_a, a_digit := nextDigit(a, ai)
		for isSpace(char_a) {
			ai++
			char_a, a_digit = nextDigit(a, ai)
		}
		char_b, b_digit := nextDigit(b, bi)
		for isSpace(char_b) {
			bi++
			char_b, b_digit = nextDigit(b, bi)
		}
		if a_digit && b_digit {
			if char_a == '0' || char_b == '0' {
				// one of these is 0 padded, so just do
				// longest-run-wins
				res := compareFractional(a, ai, b, bi)
				if res != 0 {
					return res
				}
			} else {
				res := compareDigits(a, ai, b, bi)
				if res != 0 {
					return res
				}
			}
		}
		if char_a < char_b {
			return -1
		} else if char_a > char_b {
			return 1
		}
		if char_a == 0 && char_b == 0 {
			// we're off the end of the strings, and haven't found
			// a difference yet
			return 0
		}
		ai++
		bi++
	}
	return 0
}
