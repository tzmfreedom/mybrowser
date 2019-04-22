package parser

import (
	"testing"
)

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
	if !isEOF(p) {
		t.Errorf("expect %t, actual %t", true, isEOF(p))
	}
	p = NewHtmlParser(" ")
	if isEOF(p) {
		t.Errorf("expect %t, actual %t", false, isEOF(p))
	}
	p.read()
	if !isEOF(p) {
		t.Errorf("expect %t, actual %t", true, isEOF(p))
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

func TestContainsChar(t *testing.T) {
	if !containsChar('#', "#.") {
		t.Errorf("containsChar should be true, but %t: %s", containsChar('#', "#."), "#")
	}
	if containsChar('a', "#.") {
		t.Errorf("containsChar should be false, but %t: %s", containsChar('a', "#."), "a")
	}
}
