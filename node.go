package main

import (
	"os"
)

type Node struct {
	Name string // Tag name, nil for text nodes
	Attribs map[string]string
	Content []byte // Tag content, nil for normal nodes (instead, they have a child node with the content)
	Children []*Node
	Parent *Node
}

func NewNode(parent *Node) *Node {
	sn := new(Node)
	sn.Attribs = make(map[string]string)
	sn.Children = make([]*Node, 0)
	sn.Parent = parent
	return sn
}

func (tn *Node) Args() []string {
	arr := make([]string, len(tn.Attribs))
	i := 0
	for name, val := range tn.Attribs {
		arr[i] = "-"+name+"=\""+val+"\""
		i++
	}
	return arr
}

// Creates a new child node on the current node, reallocating the slice of children and returning a pointer to the newly allocated type.
func (n *Node) NewChild() *Node {
	newNode := NewNode(n)
	n.Children = append(n.Children, newNode)
	return newNode
}

// "Writes" a string of bytes to the Content of the current tag.
func (n *Node) Write(buf []byte) (int, os.Error) {
	n.Content = append(n.Content, buf...)
	// Should never fail
	return len(buf), nil
}

// Searches up the tag heiarchy for a "normal" tag node that is a parent of this node, returning the first one found.
func (n *Node) TagParent() *Node {
	node := n
	for {
		if node.IsTag() {
			return node
		}
		node = n.Parent
		if node == nil {
			return nil
		}
	}
	panic("Not reached!")
}

// Can this node possibly have a sibling? Does it have a parent with at least two children?
func (n *Node) HasSibling() bool {
	if n.Parent == nil {
		return false
	}
	// Not == 0 because we need a second element at least to have a sibling.
	if len(n.Parent.Children) < 2 {
		return false
	}
	return true
}

func (n *Node) PrevSibling() *Node {
	if !n.HasSibling() {
		return nil
	}
	children := n.Parent.Children
	for i, node := range children {
		// Testing for POINTER equality
		if node == n {
			if i-1 >= 0 {
				return children[i-1]
			}
		}
	}
	return nil
}

func (n *Node) NextSibling() *Node {
	if !n.HasSibling() {
		return nil
	}
	children := n.Parent.Children
	for i, node := range children {
		// Testing for POINTER equality
		if node == n {
			if i+1 < len(children) {
				return children[i+1]
			}
		}
	}
	return nil
}

func (n *Node) HasChildren() bool {
	if len(n.Children) > 0 {
		return true
	}
	return false
}

// Returns the first child node furthest down the rabbit hole.
func (n *Node) EndChild() *Node {
	if !n.HasChildren() {
		return n
	}
	return n.Children[0].EndChild()
}

func (n *Node) IsTag() bool {
	return len(n.Name) != 0
}

func (n *Node) IsText() bool {
	return len(n.Content) != 0
}

func (sn *Node) String() string {
	return sn.str(" ")
}

func (sn *Node) str(indent string) string {
	if sn == nil {
		return "(nil node)"
	}
	str := ""
	if sn.IsTag() {
		str += "TAG NODE @" + sn.Name
		for key, val := range sn.Attribs {
			str += " " + key + "=" + val
		}
	} else if sn.IsText() {
		str += "TEXT NODE"
		str += " Content: '" + string(sn.Content) + "'"
	}
	for _, node := range sn.Children {
		if node == nil {
			continue
		}
		str += "\n"+indent+"+ Child: " + node.str(indent+" ")
	}
	return str
}