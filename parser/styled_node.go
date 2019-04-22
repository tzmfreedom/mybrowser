package parser

type StyledNode struct {
	Node *Node
	Properties map[string][]string
	Children []*StyledNode
}

func NewStyledNode(n *Node, styleRules []*StyleRule) *StyledNode {
	properties := map[string][]string{}
	for _, styleRule := range styleRules {
		if matchSelector(n, styleRule) {
			for _, rule := range styleRule.Rules {
				properties[rule.Key] = rule.Value
			}
		}
	}
	children := make([]*StyledNode, len(n.Children))
	for i, n := range n.Children {
		children[i] = NewStyledNode(n, styleRules)
	}
	return &StyledNode{
		Node: n,
		Properties: properties,
		Children: children,
	}
}

func matchSelector(n *Node, s *StyleRule) bool {
	for _, s := range s.Selectors {
		if n.Type != "tag" {
			continue
		}
		switch s.Identifier[0] {
		case '#':
			if n.IsEqualID(s.Identifier[1:]) {
				return true
			}
		case '.':
			if n.BelongsToClass(s.Identifier) {
				return true
			}
		default:
			if n.TagName == s.Identifier {
				return true
			}
		}
	}
	return false
}