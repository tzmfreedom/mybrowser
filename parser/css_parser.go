package parser

import (
	"errors"
	"fmt"
	"strings"
)

type StyleRule struct {
	Selectors []*Selector
	Rules     []*Rule
}

type Selector struct {
	Identifier string
}

type Rule struct {
	Key   string
	Value []string
}

type CssParser struct {
	src string
	pos int
}

func NewCssParser(src string) *CssParser {
	return &CssParser{
		src: src,
		pos: 0,
	}
}

func (p *CssParser) GetSrc() string {
	return p.src
}

func (p *CssParser) GetPos() int {
	return p.pos
}

func (p *CssParser) SetPos(pos int) {
	p.pos = pos
}

func (p *CssParser) Parse() ([]*StyleRule, error) {
	return p.parseStyles()
}

func (p *CssParser) parseStyles() ([]*StyleRule, error) {
	styles := []*StyleRule{}
	for {
		style, err := p.parseStyle()
		if err != nil {
			return styles, err
		}
		styles = append(styles, style)
		consumeSpace(p)
		if isEOF(p) {
			break
		}
	}
	return styles, nil
}

func (p *CssParser) parseStyle() (*StyleRule, error) {
	selectors, err := p.parseSelectors()
	if err != nil {
		return nil, err
	}
	rules, err := p.parseRules()
	if err != nil {
		return nil, err
	}
	return &StyleRule{
		Selectors: selectors,
		Rules:     rules,
	}, nil
}

func (p *CssParser) parseSelectors() ([]*Selector, error) {
	selectors := []*Selector{}
	for {
		consumeSpace(p)
		selector, err := p.parseSelector()
		if err != nil {
			break
		}
		selectors = append(selectors, selector)
		consumeSpace(p)
		switch peek(p) {
		case 0:
			goto L
		case '{':
			read(p)
			goto L
		case ',':
			read(p)
			continue
		default:
			return nil, errors.New("cannot parse selectors at parseSelectors")
		}
	}
L:
	return selectors, nil
}

func (p *CssParser) parseSelector() (*Selector, error) {
	c := peek(p)
	if isLetter(c) {
		l := []byte{read(p)}
		for {
			c = peek(p)
			if !isAlphaNumeric(c) {
				break
			}
			l = append(l, read(p))
		}
		return &Selector{
			Identifier: string(l),
		}, nil
	}
	if containsChar(c, "#.") {
		l := []byte{read(p)}
		c = peek(p)
		if isLetter(c) {
			for {
				c = peek(p)
				if !isAlphaNumeric(c) {
					break
				}
				l = append(l, read(p))
			}
			return &Selector{
				Identifier: string(l),
			}, nil
		}
	}
	return nil, errors.New("cannot parse selector at parseSelector")
}

func (p *CssParser) parseRules() ([]*Rule, error) {
	rules := []*Rule{}
	for {
		consumeSpace(p)
		if peek(p) == '}' {
			read(p)
			break
		}
		key, err := parseString(p)
		if err != nil {
			return nil, err
		}
		consumeSpace(p)
		if peek(p) != ':' {
			return nil, errors.New("cannot parse selector at parseRules")
		}
		read(p)
		consumeSpace(p)
		value, err := parseNotChar(p, ';')
		if err != nil {
			return nil, err
		}
		if peek(p) != ';' {
			return nil, errors.New("cannot parse selector at parseRules")
		}
		read(p)
		values := strings.Split(value, " ")

		rules = append(rules, &Rule{
			Key:   key,
			Value: values,
		})
	}
	return rules, nil
}

func (p *CssParser) parseRule() (*Rule, error) {
	return nil, nil
}

func debug() {
	fmt.Println("hoge")
}
