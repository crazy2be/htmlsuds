package main

type SudNode struct {
	Name string
	Arguments map[string]string
	Content []byte
	Children []*SudNode
	Parent *SudNode
}

func NewSudNode(parent *SudNode) *SudNode {
	sn := new(SudNode)
	sn.Arguments = make(map[string]string)
	sn.Children = make([]*SudNode, 0)
	sn.Parent = parent
	return sn
}

func (sn *SudNode) String() string {
	return sn.str(" ")
}

func (sn *SudNode) str(indent string) string {
	if sn == nil {
		return "(nil node)"
	}
	str := "NODE @" + sn.Name
	for key, val := range sn.Arguments {
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