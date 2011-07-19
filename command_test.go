package main

import (
	"testing"
)

func TestCalculateCommandPath(t *testing.T) {
	nodePath1 := []string{"form", "textbox"}
	path1 := calculateCommandPath(nodePath1)
	ex := "suds/form/textbox"
	
	if path1 != ex {
		t.Fatal("Got", path1, "expected", ex)
	}
	
	nodePath2 := []string{"somenode", "root", "foobar", "form", "textbox"}
	path2 := calculateCommandPath(nodePath2)
	
	if path2 != ex {
		t.Fatal("Got", path2, "expected", ex)
	}
}