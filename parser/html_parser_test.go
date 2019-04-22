package parser

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		input    string
		expected *Node
	}{
		{
			"<html></html>",
			&Node{
				Type:     "tag",
				TagName:  "html",
				Attr:     map[string]string{},
				Text:     "",
				Children: []*Node{},
			},
		},
		{
			"<h1></h1>",
			&Node{
				Type:     "tag",
				TagName:  "h1",
				Attr:     map[string]string{},
				Text:     "",
				Children: []*Node{},
			},
		},
		{
			"<html><head></head></html>",
			&Node{
				Type:    "tag",
				TagName: "html",
				Attr:    map[string]string{},
				Text:    "",
				Children: []*Node{
					{
						Type:     "tag",
						TagName:  "head",
						Attr:     map[string]string{},
						Text:     "",
						Children: []*Node{},
					},
				},
			},
		},
		{
			"<html><head></head><body></body></html>",
			&Node{
				Type:    "tag",
				TagName: "html",
				Attr:    map[string]string{},
				Text:    "",
				Children: []*Node{
					{
						Type:     "tag",
						TagName:  "head",
						Attr:     map[string]string{},
						Text:     "",
						Children: []*Node{},
					},
					{
						Type:     "tag",
						TagName:  "body",
						Attr:     map[string]string{},
						Text:     "",
						Children: []*Node{},
					},
				},
			},
		},
		{
			"<html foo=\"bar\" baz><head hoge=\"fuga\"></head></html>",
			&Node{
				Type:    "tag",
				TagName: "html",
				Attr: map[string]string{
					"foo": "bar",
					"baz": "",
				},
				Text: "",
				Children: []*Node{
					{
						Type:    "tag",
						TagName: "head",
						Attr: map[string]string{
							"hoge": "fuga",
						},
						Text:     "",
						Children: []*Node{},
					},
				},
			},
		},
		{
			"<div>foo<h1>bar</h1>baz</div>",
			&Node{
				Type:    "tag",
				TagName: "div",
				Attr:    map[string]string{},
				Text:    "",
				Children: []*Node{
					{
						Type:    "text",
						TagName: "",
						Text:    "foo",
					},
					{
						Type:    "tag",
						TagName: "h1",
						Attr:    map[string]string{},
						Text:    "",
						Children: []*Node{
							{
								Type:    "text",
								TagName: "",
								Text:    "bar",
							},
						},
					},
					{
						Type:    "text",
						TagName: "",
						Text:    "baz",
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		p := NewHtmlParser(testCase.input)
		actual, err := p.Parse()
		if err != nil {
			t.Errorf("raised error: %s", err.Error())
			continue
		}

		if !cmp.Equal(testCase.expected, actual) {
			fmt.Println(cmp.Diff(testCase.expected, actual))
			pp.Println(actual)
			t.Errorf("not much %s", testCase.input)
		}
	}
}

func TestParseAttr(t *testing.T) {
	testCases := []struct {
		input    string
		expected map[string]string
		error    error
	}{
		{
			input: "hoge=\"fuga\"",
			expected: map[string]string{
				"hoge": "fuga",
			},
			error: nil,
		},
		{
			input: "hoge=\"fuga\"  foo=\" bar \"",
			expected: map[string]string{
				"hoge": "fuga",
				"foo":  " bar ",
			},
			error: nil,
		},
	}
	for _, testCase := range testCases {
		p := NewHtmlParser(testCase.input)
		actual, err := p.parseAttr()
		if err != testCase.error {
			t.Errorf("expected %v, actual %v", err, testCase.error)
			continue
		}

		if !cmp.Equal(testCase.expected, actual) {
			t.Errorf("expected %s, actual %s", actual, testCase.expected)
		}
	}
}
