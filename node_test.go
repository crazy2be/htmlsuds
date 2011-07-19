package main

import (
	"testing"
	"reflect"
)

// func TestNewNode(t *testing.T) {
// 	n := NewNode(nil)
// 	//t.Print(n)
// }

func makeNodeTree() *Node {
	n0 := NewNode(nil)
	n0.Name = "fooo"
	n00 := n0.NewChild()
	n00.Name = "Foobar"
	n01 := n0.NewChild()
	n000 := n00.NewChild()
	n000.Name = "n000"
	n010 := n01.NewChild()
	n010.Name = "n010"
	return n0
}

func TestNodeTree(t *testing.T) {
	n0 := makeNodeTree()
	var ex *Node
	
	innermost := n0.EndChild()
	t.Log(innermost)
	ex = n0.FirstChild().FirstChild()
	if innermost != ex {
		t.Fatal("Innermost node calculation using EndChild() returned incorrect value! Got", innermost, ", expected", ex)
	}
	
	n01innermost := n0.FirstChild().NextSibling().EndChild()
	ex = n0.FirstChild().NextSibling().FirstChild()
	if n01innermost != ex {
		t.Fatal("Innermost node calculation using EndChild() returned incorrect value! Got", n01innermost, ", expected", ex)
	}
	
	nextinnermost := innermost.NextEndChild()
	t.Log(innermost)
	if nextinnermost != ex {
		t.Fatalf("Finding the next innermost node using NextEndChild() failed! Got %s, expected %s.\n", nextinnermost, ex)
	}
}

func TestTagPath(t *testing.T) {
	n0 := makeNodeTree()
	
	nl := n0.EndChild()
	
	tp := nl.TagPath()
	ex := []string{"fooo", "Foobar", "n000"}
	
	if !reflect.DeepEqual(tp, ex) {
		t.Fatal("When getting TagPath(), got", tp, "expected", ex)
	}
	
}

func TestTagPathLength(t *testing.T) {
	n0 := makeNodeTree()
	
	nl := n0.EndChild()
	
	pl := nl.TagPathLength()
	ex := 3
	
	t.Log(ex)
	
	if pl != ex {
		t.Fatalf("TagPathLength() returned incorrect value, got %d, expected %d.\n", pl, ex)
	}
}