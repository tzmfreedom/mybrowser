package parser

import (
	"errors"
)

type Node struct {
	Type     string
	TagName  string
	Attr     map[string]string
	Text     string
	Children []*Node
}

type HtmlParser struct {
	src string
	pos int
}

func NewHtmlParser(src string) *HtmlParser {
	return &HtmlParser{
		src: src,
		pos: 0,
	}
}

func (p *HtmlParser) Parse() (*Node, error) {
	return p.parseTag()
}

func (p *HtmlParser) peek() byte {
	if p.isEOF() {
		return 0
	}
	return p.src[p.pos]
}

func (p *HtmlParser) read() byte {
	b := p.src[p.pos]
	p.pos++
	return b
}

func (p *HtmlParser) parseTextOrTags() []*Node {
	tags := []*Node{}
	for {
		before := p.pos
		tag, err := p.parseTextOrTag()
		if err != nil || tag == nil {
			p.pos = before
			break
		}
		tags = append(tags, tag)
	}
	return tags
}

func (p *HtmlParser) parseTextOrTag() (*Node, error) {
	t, err := p.parseText()
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}
	t, err = p.parseTag()
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}
	return nil, errors.New("cannot parse at #parseTextOrTag")
}

func (p *HtmlParser) parseText() (*Node, error) {
	text, err := p.parseNotString("<>")
	if err != nil {
		return nil, err
	}
	if text == "" {
		return nil, nil
	}
	return &Node{
		Type:    "text",
		TagName: "",
		Text:    text,
	}, nil
}

func (p *HtmlParser) parseTag() (*Node, error) {
	if p.peek() != '<' {
		return nil, errors.New("not parse")
	}
	p.read()

	tagName, err := p.parseString()
	if err != nil {
		return nil, err
	}
	attr, err := p.parseAttr()
	if p.peek() != '>' {
		return nil, errors.New("not parse")
	}
	p.read()

	children := p.parseTextOrTags()
	if !p.readString("</") {
		return nil, errors.New("cannot parse")
	}

	endTagName, err := p.parseString()
	if err != nil {
		return nil, err
	}
	if endTagName != tagName {
		return nil, errors.New("cannot much start/end tagname")
	}
	if p.peek() != '>' {
		return nil, errors.New("not parse")
	}
	p.read()

	return &Node{
		Type:     "tag",
		TagName:  tagName,
		Attr:     attr,
		Text:     "",
		Children: children,
	}, nil
}

func (p *HtmlParser) parseAttr() (map[string]string, error) {
	attr := map[string]string{}
	for {
		p.consumeSpace()
		attrKey, err := p.parseString()
		if err != nil {
			return attr, nil
		}
		if p.peek() != '=' {
			attr[attrKey] = ""
			continue
		}
		p.read()
		if p.peek() != '"' {
			return map[string]string{}, errors.New("cannot parse attribute")
		}
		p.read()
		attrValue, err := p.parseNotChar('"')
		if err != nil {
			return map[string]string{}, err
		}
		attr[attrKey] = attrValue
		p.read()
	}
}

func (p *HtmlParser) parseString() (string, error) {
	l := p.peek()
	if !isLetter(l) {
		return "", errors.New("can not parse string")
	}
	ret := []byte{l}
	p.read()
	for {
		c := p.peek()
		if isSpace(c) {
			break
		}
		if !isAlphaNumeric(c) {
			break
		}
		p.read()
		ret = append(ret, c)
	}
	return string(ret), nil
}

func (p *HtmlParser) parseNotChar(b byte) (string, error) {
	l := []byte{}
	for {
		c := p.peek()
		if c == b {
			return string(l), nil
		}
		p.read()
		l = append(l, c)
	}
}

func (p *HtmlParser) parseNotString(src string) (string, error) {
	l := []byte{}
	for {
		c := p.peek()
		if containsChar(c, src) {
			return string(l), nil
		}
		p.read()
		l = append(l, c)
	}
}

func (p *HtmlParser) consumeSpace() {
	for p.peek() == ' ' {
		p.read()
	}
}

func (p *HtmlParser) readString(src string) bool {
	for i, s := range src {
		if p.src[p.pos+i] != byte(s) {
			return false
		}
	}
	p.pos = p.pos + len(src)
	return true
}

func (p *HtmlParser) isEOF() bool {
	return len(p.src) <= p.pos
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
	return b == ' '
}

func containsChar(b byte, src string) bool {
	for _, s := range src {
		if byte(s) == b {
			return true
		}
	}
	return false
}
