package gdefer

import "strings"

type node struct {
	/*
		以/p/:lang/doc為例
		pattern 指的是 /p/:lang
		part 指的是 :lang
		children 指的是 /doc
		isWild 指的是是否存在 * 或 : ->有的話true
	*/
	pattern  string  //帶匹配路由
	part     string  //路由中的一部分
	children []*node //子節點
	isWild   bool    //是否精確匹配
}

//查詢第一個匹配的節點，用於插入
func (n *node) matchChildren(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配節點，用於查找
func (n *node) allMatchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	children := n.matchChildren(part)
	if children == nil {
		children = &node{
			part:   part,
			isWild: part[0] == '*' || part[0] == ':',
		}
		n.children = append(n.children, children)
	}
	children.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.allMatchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
