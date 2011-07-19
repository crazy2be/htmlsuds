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
	Processed bool // Has this node been processed yet?
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

func (n *Node) FirstChild() *Node {
	if !n.HasChildren() {
		return nil
	}
	return n.Children[0]
}

// Returns the first child node furthest down the rabbit hole.
func (n *Node) EndChild() *Node {
	if !n.HasChildren() {
		if n.IsTag() {
			return n
		}
		return n.Parent
	}
	return n.Children[0].EndChild()
}

// Returns the next end child from this node. If this is not the end node in the node tree, then this will return the same thing as EndChild(). Otherwise, it returns the EndChild() of the next sibling or parent's sibling or parent's parent's sibling, ad infinitum.
func (n *Node) NextEndChild() *Node {
	//fmt.Println("In NextEndChild():", n)
// 	ec := n.EndChild()
// 	if ec != n {
// 		//fmt.Println(ec, n)
// 		return ec
// 	}
	if n.HasSibling() {
		ns := n.NextSibling()
		if ns != nil {
			return ns.EndChild()
		}
	}
	if n.Parent != nil {
		ps := n.Parent.NextSibling()
		if ps != nil {
			return ps.NextEndChild()
		}
	}
	return nil
}

// Returns the path to the tag as a series of strings. That is, {"rootName", "parentname", "nodename"} for the following:
// @parentname
//  @nodename
func (n *Node) TagPath() []string {
	i := n.TagPathLength()
	tp := make([]string, i)
	currentnode := new(int)
	*currentnode = i-1
	n.tagPath(tp, currentnode)
	return tp
}

func (n *Node) tagPath(paths []string, currentnode *int) {
	paths[*currentnode] = n.Name
	*currentnode--
	if n.Parent != nil {
		n.Parent.tagPath(paths, currentnode)
	}
}

func (n *Node) TagPathLength() int {
	i := new(int)
	n.tagPathLength(i)
	return *i
}

func (n *Node) tagPathLength(i *int) {
	*i++
	if n.Parent != nil {
		n.Parent.tagPathLength(i)
		return
	}
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
	} else if sn.IsText() {
		str += "TEXT NODE"
	} else {
		str += "UNKNOWN NODE"
	}
	if sn.IsTag() {
		for key, val := range sn.Attribs {
			str += " " + key + "=" + val
		}
	}
	if sn.IsText() {
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