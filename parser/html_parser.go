package parser

import (
	"errors"
	"strings"
)

type Node struct {
	Type     string
	TagName  string
	Attr     map[string]string
	Text     string
	Children []*Node
}

func (n *Node) BelongsToClass(class string) bool {
	if _, ok := n.Attr["class"]; ok {
		classes := strings.Split(strings.TrimSpace(class), " ")
		for _, c := range classes {
			if c == class {
				return true
			}
		}
	}
	return false
}

func (n *Node) IsEqualID(id string) bool {
	v, ok := n.Attr["id"]
	return ok && v == id
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

func (p *HtmlParser) GetSrc() string {
	return p.src
}

func (p *HtmlParser) GetPos() int {
	return p.pos
}

func (p *HtmlParser) SetPos(pos int) {
	p.pos = pos
}

func (p *HtmlParser) Parse() (*Node, error) {
	return p.parseTag()
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
	text, err := parseNotString(p, "<>")
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
	if peek(p) != '<' {
		return nil, errors.New("not parse")
	}
	read(p)

	tagName, err := parseString(p)
	if err != nil {
		return nil, err
	}
	attr, err := p.parseAttr()
	if peek(p) != '>' {
		return nil, errors.New("not parse")
	}
	read(p)

	children := p.parseTextOrTags()
	if !readString(p, "</") {
		return nil, errors.New("cannot parse")
	}

	endTagName, err := parseString(p)
	if err != nil {
		return nil, err
	}
	if endTagName != tagName {
		return nil, errors.New("cannot much start/end tagname")
	}
	if peek(p) != '>' {
		return nil, errors.New("not parse")
	}
	read(p)

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
		consumeSpace(p)
		attrKey, err := parseString(p)
		if err != nil {
			return attr, nil
		}
		if peek(p) != '=' {
			attr[attrKey] = ""
			continue
		}
		read(p)
		if peek(p) != '"' {
			return map[string]string{}, errors.New("cannot parse attribute")
		}
		read(p)
		attrValue, err := parseNotChar(p, '"')
		if err != nil {
			return map[string]string{}, err
		}
		attr[attrKey] = attrValue
		read(p)
	}
}
