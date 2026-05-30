package magc

import (
	"errors"
	"github.com/addyreal/goutils/magc/internal/rules"
	"sort"
)

type b struct {
	rules []*rules.Rule
}

type Table struct {
	offsets []int
	anchors map[int]map[int][]*rules.Rule
}

func Init() *b {
	return &b{}
}

func (x *b) Add(res any, pattern ...int) error {
	if len(pattern) == 0 {
		return nil
	}

	for _, e := range pattern {
		if e < -1 || e > 255 {
			return errors.New("byte value neither -1 or 0-255")
		}
	}

	rule, err := rules.MakeRule(res, pattern)
	if err != nil {
		return err
	}

	x.rules = append(x.rules, rule)

	return nil
}

func (x *b) AddMust(res any, pattern ...int) {
	err := x.Add(res, pattern...)
	if err != nil {
		panic(err.Error())
	}
}

func (x *b) Get() *Table {
	res := &Table{
		anchors: make(map[int]map[int][]*rules.Rule),
	}

	rules.SortRules(x.rules)

	for _, r := range x.rules {
		offset, anchor := r.GetAnchor()
		if res.anchors[offset] == nil {
			res.anchors[offset] = make(map[int][]*rules.Rule)
		}

		res.offsets = append(res.offsets, offset)
		res.anchors[offset][anchor] = append(
			res.anchors[offset][anchor],
			r,
		)
	}

	sort.Ints(res.offsets)

	return res
}

func (t *Table) Find(b []byte) (any, bool) {
	if len(b) == 0 {
		return nil, false
	}

	for _, offset := range t.offsets {
		if len(b) <= offset {
			continue
		}

		v := int(b[offset])
		candidates, ok := t.anchors[offset][v]
		if ok {
			for _, c := range candidates {
				if c.MatchesBytes(b) {
					return c.GetResult(), true
				}
			}
		}
	}

	return nil, false
}
