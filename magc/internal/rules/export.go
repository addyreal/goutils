package rules

import (
	"errors"
	"sort"
)

type Rule struct {
	offset  int
	pattern []int
	result  any
}

func MakeRule(res any, pattern []int) (*Rule, error) {
	if len(pattern) == 0 {
		return nil, errors.New("pattern is empty")
	}

	for i, v := range pattern {
		if v != -1 {
			return &Rule{
				offset:  i,
				result:  res,
				pattern: pattern,
			}, nil
		}
	}

	return nil, errors.New("pattern is only wildcards")
}

// Sorts rules longest to shortest to prioritize verbose rules
func SortRules(rules []*Rule) {
	sort.Slice(rules, func(i, j int) bool {
		return len(rules[i].pattern) > len(rules[j].pattern)
	})
}

func (r *Rule) GetAnchor() (offset int, anchor int) {
	return r.offset, r.pattern[r.offset]
}

func (r *Rule) GetResult() any {
	return r.result
}

// Checks whether bytes match a pattern
func (r *Rule) MatchesBytes(b []byte) bool {
	if len(b) <= r.offset {
		return false
	}

	if len(b) < len(r.pattern) {
		return false
	}

	for i, x := range r.pattern {
		if x == -1 {
			continue
		}

		if b[i] != byte(x) {
			return false
		}
	}

	return true
}
