package parser

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"

	"github.com/google/go-cmp/cmp"
)

func TestParseCss(t *testing.T) {
	testCases := []struct {
		input    string
		expected []*StyleRule
	}{
		{
			`#foo{}`,
			[]*StyleRule{
				{
					Selectors: []*Selector{
						{
							Identifier: "#foo",
						},
					},
					Rules: []*Rule{},
				},
			},
		},
		{
			`.foo{}`,
			[]*StyleRule{
				{
					Selectors: []*Selector{
						{
							Identifier: ".foo",
						},
					},
					Rules: []*Rule{},
				},
			},
		},
		{
			`div{}`,
			[]*StyleRule{
				{
					Selectors: []*Selector{
						{
							Identifier: "div",
						},
					},
					Rules: []*Rule{},
				},
			},
		},
		{
			`.foo, #bar, div{}`,
			[]*StyleRule{
				{
					Selectors: []*Selector{
						{
							Identifier: ".foo",
						},
						{
							Identifier: "#bar",
						},
						{
							Identifier: "div",
						},
					},
					Rules: []*Rule{},
				},
			},
		},
		{
			`div{
display: block;
margin: 0 10px 2px 0;
}`,
			[]*StyleRule{
				{
					Selectors: []*Selector{
						{
							Identifier: "div",
						},
					},
					Rules: []*Rule{
						{
							Key:   "display",
							Value: []string{"block"},
						},
						{
							Key:   "margin",
							Value: []string{"0", "10px", "2px", "0"},
						},
					},
				},
			},
		},
	}
	for i, testCase := range testCases {
		p := NewCssParser(testCase.input)
		actual, err := p.Parse()
		if err != nil {
			pp.Println(i)
			t.Errorf("raised error: %s", err.Error())
			continue
		}

		if !cmp.Equal(testCase.expected, actual) {
			fmt.Println(cmp.Diff(testCase.expected, actual))
			pp.Println(i)
			pp.Println(actual)
			t.Errorf("not much %s", testCase.input)
		}
	}
}
