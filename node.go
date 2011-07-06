package main

type Node struct {
	Name string
	Args map[string]string
	Content []byte
	Children []*Node
	Parent *Node
}

func NewNode(parent *Node) *Node {
	sn := new(Node)
	sn.Args = make(map[string]string)
	sn.Children = make([]*Node, 0)
	sn.Parent = parent
	return sn
}

func NewSudNode(parent *Node) *Node {
	return NewNode(parent)
}

func (sn *Node) String() string {
	return sn.str(" ")
}

func (sn *Node) str(indent string) string {
	if sn == nil {
		return "(nil node)"
	}
	str := "NODE @" + sn.Name
	for key, val := range sn.Args {
		str += " " + key + "=" + val
	}
	if len(sn.Content) > 0 {
		str += " [Content: '" + string(sn.Content) + "']"
	}
	for _, node := range sn.Children {
		if node == nil {
			continue
		}
		str += "\n"+indent+"+ Child:" + node.str(indent+" ")
	}
	return str
}