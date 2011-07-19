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
		return nodePath[0]
	}
	return calculateCommandPath(nodePath[1:])
}

func (n *Node) CalculateEnv() []string {
	curenv := os.Environ()
	env := make([]string, len(n.Attribs) + len(curenv))
	i := 0
	for _, _ = range curenv {
		env[i] = curenv[i]
		i++
	}
	for key, val := range n.Attribs {
		env[i] = "ARGS_" + strings.ToUpper(key) + "=" + val
		i++
	}
	return env
}

func processNode(node *Node) (err os.Error) {
	if node == nil {
		err = os.NewError("Nil node passed, cannot possibly generate HTML.")
		return
	}
	
	cmdname := calculateCommandPath(node.TagPath())
	log.Println("For node:", node.Name)
	log.Println("Calculated command name as:", cmdname)
	cmdline := node.Args()
	log.Println("Calculated commandline as:", cmdline)
	cmd := exec.Command(cmdname, cmdline...)
	
	cmd.Env = node.CalculateEnv()
	cmd.Stderr = os.Stderr
	input, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	
	//input := bytes.NewBuffer([]byte(""))
	go func() {
		input.Write(node.Content)
		input.Close()
	}()
	//go input.WriteTo(inputp)
	
	node.Content, err = cmd.Output()
	if err != nil {
		return
	}
	log.Println(node)
	
	node.Processed = true
	
	return
}