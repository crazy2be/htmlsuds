package main

import (
	"bytes"
	"exec"
	"log"
	"os"
	"io"
)

func GenerateHTML(root *Node, wr io.Writer) (os.Error) {
	val, err := generateHTML(root, wr)
	wr.Write(val)
	return err
}

func generateHTML(root *Node, wr io.Writer) (val []byte, err os.Error) {
	if root == nil {
		err = os.NewError("Nil node passed, cannot possibly generate HTML.")
		return
	}
	cmdname := root.Name
	cmdline := root.Args()
	cmd := exec.Command("suds/"+cmdname, cmdline...)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	inputp, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	input := bytes.NewBuffer([]byte(""))
	input.Write(root.Content)
	go input.WriteTo(inputp)
	buf := bytes.NewBuffer([]byte(""))
	buf.ReadFrom(output)
	raw := buf.Bytes()
	subroot, err := Parse(buf)
	if err != nil {
		return
	}
	if len(subroot.Children) == 0 {
		wr.Write(raw)
		return
	}
	GenerateHTML(subroot, wr)
	cmd.Run()
	log.Println(cmd)
	return
}