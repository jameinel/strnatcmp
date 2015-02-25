// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"sort"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/testing"
)

type naturalSortSuite struct {
	testing.BaseSuite
}

var _ = gc.Suite(&naturalSortSuite{})

func (s *naturalSortSuite) TestNaturallyEmpty(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{},
		[]string{},
	)
}

func (s *naturalSortSuite) TestNaturallyAlpha(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"bac", "cba", "abc"},
		[]string{"abc", "bac", "cba"},
	)
}

func (s *naturalSortSuite) TestNaturallyAlphaNumeric(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"a1", "a10", "a100", "a11"},
		[]string{"a1", "a10", "a11", "a100"},
	)
}

func (s *naturalSortSuite) TestNaturallySpecial(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"a1", "a10", "a100", "a1/1", "1a"},
		[]string{"1a", "a1", "a1/1", "a10", "a100"},
	)
}

func (s *naturalSortSuite) TestNaturallyTagLike(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"a1/1", "a1/11", "a1/2", "a1/7", "a1/100"},
		[]string{"a1/1", "a1/2", "a1/7", "a1/11", "a1/100"},
	)
}

func (s *naturalSortSuite) TestNaturallyMixed(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"a1a", "a", "a1", "a1b", "a0", "a20", "a2", "a10"},
		[]string{"a", "a0", "a1", "a1a", "a1b", "a2", "a10", "a20"},
	)
}

func (s *naturalSortSuite) TestNaturallySeveralNumericParts(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"x2-y08", "x2-g8", "x8-y8", "x2-y7"},
		[]string{"x2-g8", "x2-y7", "x2-y08", "x8-y8"},
	)
}

func (s *naturalSortSuite) TestNaturallyDecimalFractions(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"1.002", "1.3", "1.001", "1.010", "1.02", "1.1"},
		[]string{"1.001", "1.002", "1.010", "1.02", "1.1", "1.3"},
	)
}

func (s *naturalSortSuite) TestNaturallyIPs(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{"100.001.010.123", "001.001.010.123", "001.002.010.123"},
		[]string{"001.001.010.123", "001.002.010.123", "100.001.010.123"},
	)
}

func (s *naturalSortSuite) TestNaturallyJuju(c *gc.C) {
	s.assertNaturallySort(
		c,
		[]string{
			"ubuntu/0",
			"ubuntu/1",
			"ubuntu/10",
			"ubuntu/100",
			"ubuntu/101",
			"ubuntu/102",
			"ubuntu/103",
			"ubuntu/104",
			"ubuntu/11"},
		[]string{
			"ubuntu/0",
			"ubuntu/1",
			"ubuntu/10",
			"ubuntu/11",
			"ubuntu/100",
			"ubuntu/101",
			"ubuntu/102",
			"ubuntu/103",
			"ubuntu/104"},
	)
}

func (s *naturalSortSuite) assertNaturallySort(c *gc.C, sample, expected []string) {
	sort.Sort(naturally(sample))
	c.Assert(sample, gc.DeepEquals, expected)
}

func makeLongNumberString(sep string) string {
	parts := []string{"start"}
	for i := 0; i < 1000; i++ {
		parts = append(parts, fmt.Sprintf("%d", i))
	}
	parts = append(parts, ".")
	for i := 0; i < 1000; i++ {
		parts = append(parts, fmt.Sprintf("%d", i+500))
	}
	return strings.Join(parts, sep)
}

func (s *naturalSortSuite) BenchmarkCompareLongStrings(c *gc.C) {
	start := makeLongNumberString("")
	longA := start + "a"
	longB := start + "b"
	for i := 0; i < c.N; i++ {
		sort.Sort(naturally([]string{longA, longB}))
	}
}

func (s *naturalSortSuite) BenchmarkLotsOfParts(c *gc.C) {
	start := makeLongNumberString(".")
	longA := start + "a"
	longB := start + "b"
	for i := 0; i < c.N; i++ {
		sort.Sort(naturally([]string{longA, longB}))
	}
}
