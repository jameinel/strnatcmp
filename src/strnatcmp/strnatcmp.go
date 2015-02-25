package strnatcmp

// isSpace is the Latin subset of unicode.IsSpace
func isSpace(c uint8) bool {
	// Original function used a switch statement, but switch can't be
	// inlined in go-1.2
	return (c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\v' ||
		c == '\r' ||
		c == 0x85 ||
		c == 0xA0)
}

// isDigit is the same as unicode.IsDigit() for Latin characters.
func isDigit(c uint8) bool {
	return '0' <= c && c <= '9'
}

// nextDigit returns the value of the next character, and whether it is a digit
func nextDigit(s string, i int) (uint8, bool) {
	if i < len(s) {
		c := s[i]
		// we only handle latin digits, but this gets tricky if you
		// start considering non-latin numbers.
		return c, isDigit(c)
	}
	return 0, false
}

// nextDigitNoSpace returns the value of the next character ignoring spaces. It
// also increments the counter if it has to skip a space character.
func nextDigitNoSpace(s string, i *int) (uint8, bool) {
	var c uint8
	for *i < len(s) {
		c = s[*i]
		if !isSpace(c) {
			break
		}
		(*i)++
	}
	return c, isDigit(c)
}

// compareNumber considers the current text as a number, and compares them.
// The longest run of digits wins, as it is clearly larger. However, if we get
// the same number of characters, then we have to consider the actual values.
// As we walk the characters we notice the first time they don't match, and
// which side has the larger value.
// This returns the comparison, and how many characters were evaluated if they are equal.
func compareNumber(a string, ai *int, b string, bi *int) int {
	var bias int = 0
	// ignore leading zeros if we aren't in a fraction
	char_a, a_digit := nextDigit(a, *ai)
	for char_a == '0' {
		(*ai)++
		char_a, a_digit = nextDigit(a, *ai)
	}
	char_b, b_digit := nextDigit(b, *bi)
	for char_b == '0' {
		(*bi)++
		char_b, b_digit = nextDigit(b, *bi)
	}
	for {
		// We don't need to check if ai or bi is off the end of the
		// string, because one of them won't be a digit, which
		// guarantees a return from this function.
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
		(*ai)++
		char_a, a_digit = nextDigit(a, *ai)
		(*bi)++
		char_b, b_digit = nextDigit(b, *bi)
	}
}

// compareFractional handles when numbers start with a 0, to handle fractional
// decimals.  (1.01 is longer than 1.2, but comes before 1.2).
// returns the comparison, and the number of characters evaluated
func compareFractional(a string, ai int, b string, bi int) (int, int) {
	// We just treat the first value to be larger as winning
	start := ai
	for {
		// We don't need to check if ai or bi is off the end of the
		// string, because one of them won't be a digit, which
		// guarantees a return from this function.
		char_a, a_digit := nextDigit(a, ai)
		char_b, b_digit := nextDigit(b, bi)
		if !a_digit && !b_digit {
			// Same length, and no differences
			return 0, ai - start
		} else if !a_digit {
			// no more digits in a, b has at least 1 extra digit
			return -1, ai - start
		} else if !b_digit {
			// no more digits in b, a has at least 1 extra digit
			return 1, ai - start
		} else if char_a < char_b {
			return -1, ai - start
		} else if char_a > char_b {
			return 1, ai - start
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
	fractional := false
	ai := 0
	bi := 0
	for {
		char_a, a_digit := nextDigitNoSpace(a, &ai)
		char_b, b_digit := nextDigitNoSpace(b, &bi)
		if a_digit && b_digit {
			if fractional {
				res, count := compareFractional(a, ai, b, bi)
				if res != 0 {
					return res
				}
				fractional = false
				ai += count
				bi += count
				continue
			} else {
				res := compareNumber(a, &ai, b, &bi)
				if res != 0 {
					return res
				}
				continue
			}
		}
		if char_a == '.' && char_b == '.' {
			fractional = true
		} else {
			fractional = false
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
}
