package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
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
			"<html foo=\"bar\"><head hoge=\"fuga\"></head></html>",
			&Node{
				Type:    "tag",
				TagName: "html",
				Attr: map[string]string{
					"foo": "bar",
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
	}
	for _, testCase := range testCases {
		p := NewHtmlParser(testCase.input)
		actual, err := p.Parse()
		if err != nil {
			t.Errorf("raised error: %s", err.Error())
			continue
		}

		if !cmp.Equal(testCase.expected, actual) {
			//fmt.Println(cmp.Diff(testCase.expected, actual))
			pp.Println(actual)
			t.Errorf("not much %s", testCase.input)
		}
	}
}

func TestPeek(t *testing.T) {
	src := "foo"
	p := NewHtmlParser(src)
	actual := p.peek()
	if actual != 'f' {
		t.Errorf("expected f, actual %s", string(actual))
	}
	if p.pos != 0 {
		t.Errorf("expected 0, actual %d", p.pos)
	}
}

func TestParseNotChar(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"\"",
			"",
		},
		{
			"123abc$\"",
			"123abc$",
		},
	}
	for _, testCase := range testCases {
		p := NewHtmlParser(testCase.input)
		actual, err := p.parseNotChar('"')
		if err != nil {
			t.Errorf("raise error: %s", err.Error())
		}
		if testCase.expected != actual {
			t.Errorf("expect %s, actual %s", testCase.expected, actual)
		}
	}
}

func TestConsumeSpace(t *testing.T) {
	p := NewHtmlParser("   ")
	p.consumeSpace()
	if p.pos != 3 {
		t.Errorf("expect %d, actual %d", 3, p.pos)
	}
}

func TestIsEOF(t *testing.T) {
	p := NewHtmlParser("")
	if !p.isEOF() {
		t.Errorf("expect %t, actual %t", true, p.isEOF())
	}
	p = NewHtmlParser(" ")
	if p.isEOF() {
		t.Errorf("expect %t, actual %t", false, p.isEOF())
	}
	p.read()
	if !p.isEOF() {
		t.Errorf("expect %t, actual %t", true, p.isEOF())
	}
}

func TestReadString(t *testing.T) {
	p := NewHtmlParser("</hoge>")
	actual := p.readString("</")
	if !actual {
		t.Errorf("expect %t, actual %t", true, actual)
	}
	if p.pos != 2 {
		t.Errorf("expect %d, actual %d", 2, p.pos)
	}
	p = NewHtmlParser("//hoge>")
	actual = p.readString("</")
	if actual {
		t.Errorf("expect %t, actual %t", false, actual)
	}
	if p.pos != 0 {
		t.Errorf("expect %d, actual %d", 0, p.pos)
	}
	p = NewHtmlParser("<<hoge>")
	actual = p.readString("</")
	if actual {
		t.Errorf("expect %t, actual %t", false, actual)
	}
	if p.pos != 0 {
		t.Errorf("expect %d, actual %d", 0, p.pos)
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

func TestIsLetter(t *testing.T) {
	success := "abcdefghijklmnopqrstuvwxyz"
	for _, s := range success {
		if !isLetter(byte(s)) {
			t.Errorf("isLetter should be true, but %t: %s", isLetter(byte(s)), string(s))
		}
	}
	failure := "0123456789"
	for _, s := range failure {
		if isLetter(byte(s)) {
			t.Errorf("isLetter should be false, but %t: %s", isLetter(byte(s)), string(s))
		}
	}
}

func TestIsNumeric(t *testing.T) {
	success := "abcdefghijklmnopqrstuvwxyz"
	for _, s := range success {
		if isNumeric(byte(s)) {
			t.Errorf("isNumeric should be false, but %t: %s", isNumeric(byte(s)), string(s))
		}
	}
	failure := "0123456789"
	for _, s := range failure {
		if !isNumeric(byte(s)) {
			t.Errorf("isNumeric should be true, but %t: %s", isNumeric(byte(s)), string(s))
		}
	}
}
