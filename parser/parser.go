package parser

import (
	"errors"
)

type Parser interface {
	GetSrc() string
	GetPos() int
	SetPos(pos int)
}

func next(p Parser) byte {
	p.SetPos(p.GetPos() + 1)
	return p.GetSrc()[p.GetPos()]
}

func peek(p Parser) byte {
	if isEOF(p) {
		return 0
	}
	return p.GetSrc()[p.GetPos()]
}

func read(p Parser) byte {
	b := p.GetSrc()[p.GetPos()]
	p.SetPos(p.GetPos() + 1)
	return b
}

func parseString(p Parser) (string, error) {
	l := peek(p)
	if !isLetter(l) {
		return "", errors.New("can not parse string")
	}
	ret := []byte{l}
	read(p)
	for {
		c := peek(p)
		if isSpace(c) {
			break
		}
		if !isAlphaNumeric(c) {
			break
		}
		read(p)
		ret = append(ret, c)
	}
	return string(ret), nil
}

func parseNotChar(p Parser, b byte) (string, error) {
	l := []byte{}
	for {
		c := peek(p)
		if c == b {
			return string(l), nil
		}
		read(p)
		l = append(l, c)
	}
}

func parseNotString(p Parser, src string) (string, error) {
	l := []byte{}
	for {
		c := peek(p)
		if containsChar(c, src) {
			return string(l), nil
		}
		read(p)
		l = append(l, c)
	}
}

func consumeSpace(p Parser) {
	for isSpace(peek(p)) {
		read(p)
	}
}

func readString(p Parser, src string) bool {
	for i, s := range src {
		if p.GetSrc()[p.GetPos()+i] != byte(s) {
			return false
		}
	}
	p.SetPos(p.GetPos() + len(src))
	return true
}

func isEOF(p Parser) bool {
	return len(p.GetSrc()) <= p.GetPos()
}

func isAlphaNumeric(b byte) bool {
	return isLetter(b) || isNumeric(b)
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t'
}

func containsChar(b byte, src string) bool {
	for _, s := range src {
		if byte(s) == b {
			return true
		}
	}
	return false
}
