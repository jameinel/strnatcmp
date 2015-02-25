package strnatcmp

import (
	"testing"

	gc "gopkg.in/check.v1"
)

func TestAll(t *testing.T) {
	gc.TestingT(t)
}

type StrNatCmpSuite struct{}

var _ = gc.Suite(&StrNatCmpSuite{})

func (s *StrNatCmpSuite) checkLess(c *gc.C, a, b string) {
	c.Check(Compare(a, b), gc.Equals, -1, gc.Commentf("%q should be less than %q", a, b))
	c.Check(Compare(b, a), gc.Equals, 1, gc.Commentf("%q should be greater than %q", b, a))
}

func (s *StrNatCmpSuite) checkSame(c *gc.C, a, b string) {
	c.Check(Compare(a, b), gc.Equals, 0, gc.Commentf("%q should be the same as %q", a, b))
	c.Check(Compare(b, a), gc.Equals, 0, gc.Commentf("%q should be the same as %q", a, b))
}

func (s *StrNatCmpSuite) TestSimple(c *gc.C) {
	// Test just that normal string sorting works
	s.checkSame(c, "", "")
	s.checkLess(c, "a", "b")
	s.checkLess(c, "", "a")
	s.checkLess(c, "a", "ab")
	s.checkLess(c, "ab", "ac")
	s.checkSame(c, "a", "a")
	s.checkSame(c, "abcd", "abcd")
}

func (s *StrNatCmpSuite) TestCompareWithInt(c *gc.C) {
	// Test that strings including numbers treat the numbers specially
	s.checkLess(c, "a1", "a2")
	s.checkLess(c, "a2", "a10")
	s.checkLess(c, "a1", "a11")
	s.checkLess(c, "a111", "a121")
	s.checkLess(c, "a121", "a1111")
	s.checkLess(c, "a121", "a1131")
	s.checkLess(c, "a113", "a121")
}

func (s *StrNatCmpSuite) TestWithLeadingZeros(c *gc.C) {
	// Test that numbers with padding zeros are handled correctly
	s.checkLess(c, "a01", "a02")
	s.checkLess(c, "a01", "a2")
	s.checkLess(c, "a02", "a1")
	s.checkLess(c, "a1.2", "a1.3")
	s.checkLess(c, "a1.02", "a1.3")
	// I think it actually gets this wrong
	s.checkLess(c, "a1.04", "a1.30")
}

func (s *StrNatCmpSuite) TestIgnoreSpaces(c *gc.C) {
	s.checkSame(c, "a1", "a 1")
	s.checkSame(c, "a1", "a\t1")
	s.checkSame(c, "a1", "a\r1")
	s.checkSame(c, "a1", "a\n1")
	s.checkSame(c, "a1", "a\v1")
	s.checkLess(c, "a1", "a 2")
	s.checkLess(c, "a2", "a 10")
	s.checkLess(c, "a b", "a   c")
	s.checkSame(c, "a b", "a   b")
	s.checkLess(c, "a b", "a   bc")
	s.checkLess(c, "a    b", "a bc")
}

// from sourcefrog.net/projects/natsort/example-out.txt
// note that the *actual* ordering in example-out.txt is from a different
// version of the sorting algorithm
var corpus = []string{
	"1-02",
	"1-2",
	"1-20",
	"10-20",
	"fred",
	"jane",
	"pic01",
	"pic02",
	"pic02a",
	"pic02000",
	"pic05",
	"pic2",
	"pic3",
	"pic4",
	"pic 4 else",
	"pic 5",
	// Original corpus has "pic 5 " but since whitespace is ignored, that
	// means it is actually equal to "pic 5" which means, we can't say
	// which one comes first
	//"pic 5 ",
	"pic 5 something",
	"pic 6",
	"pic   7",
	"pic100",
	"pic100a",
	"pic120",
	"pic121",
	"tom",
	"x2-g8",
	"x2-y08",
	"x2-y7",
	"x8-y8",
}

func (s *StrNatCmpSuite) TestCorpus(c *gc.C) {
	// test that all strings in the corpus sort relative to all other
	// strings
	for i := range corpus {
		s.checkSame(c, corpus[i], corpus[i])
	}
	for i, a := range corpus {
		for _, b := range corpus[i+1:] {
			s.checkLess(c, a, b)
		}
	}
}
