package main

import (
	"strings"
	"exec"
	"log"
	"os"
	"io"
	"github.com/crazy2be/osutil"
)

func GenerateHTML(root *Node, wr io.Writer) (err os.Error) {
	node := root.EndChild()
	for {
		err = processNode(node)
		if err != nil {
			return err
		}
		node = node.NextEndChild()
		if node == nil {
			log.Println("Note: Writing all processed nodes to writer not yet implemented! Here's what the layout looks like:\n", root)
			return
		}
	}
	return
}

// Looks through all the possible command paths for the given node, and resolves to the most specific one.
func calculateCommandPath(nodePath []string) string {
	path := "suds/" + strings.Join(nodePath, "/") + "/main"
	if osutil.FileExists(path) {
		return path
	}
	if len(nodePath) < 2 {
		return ""
	}
	return calculateCommandPath(nodePath[1:])
}

func processNode(node *Node) (err os.Error) {
	if node == nil {
		err = os.NewError("Nil node passed, cannot possibly generate HTML.")
		return
	}
	
	cmdname := calculateCommandPath(node.TagPath())
	log.Println("Calculated command name as:", cmdname)
	cmdline := node.Args()
	log.Println("Calculated commandline as:", cmdline)
	cmd := exec.Command(cmdname, cmdline...)
	
	cmd.Stderr = os.Stderr
	input, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	
	//input := bytes.NewBuffer([]byte(""))
	go input.Write(node.Content)
	//go input.WriteTo(inputp)
	
	node.Content, err = cmd.Output()
	if err != nil {
		return
	}
	log.Println(node)
	
	node.Processed = true
	
	return
}