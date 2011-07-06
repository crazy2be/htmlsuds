package main

import (
	"exec"
	"log"
	"os"
	"io"
)

func GenerateHTML(root *Node, wr io.Writer) (os.Error) {
	if root == nil {
		return os.NewError("Nil node passed, cannot possibly generate HTML.")
	}
	cmdname := root.Name
	cmdline := make([]string, len(root.Args))
	i := 0
	for key, val := range root.Args {
		cmdline[i] = "-"+key+"=\""+val+"\""
		i++
	}
	cmd := exec.Command("suds/"+cmdname, cmdline...)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	input, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	cmd.Run()
	log.Println(cmd)
	return nil
}