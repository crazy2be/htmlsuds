package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"bufio"
)

type parseState struct {
	inTag bool
	inQuotes bool
	inTagBody bool
	tagClosed bool
	lineNo int
	currentNode *Node
	rootNode *Node
	rd io.Reader
	wr io.Writer
}

func newParseState(rd io.Reader, wr io.Writer) *parseState {
	ps := new(parseState)
	ps.rootNode = NewNode(nil)
	ps.currentNode = ps.rootNode
	ps.rootNode.Name = "root"
	log.Println(ps)
	ps.rd = rd
	ps.wr = wr
	ps.lineNo = -1
	return ps
}

func (p *parseState) parseLine(l []uint8) {
	p.lineNo++
	if len(l) < 1 {
		return
	}
	log.Println("In parseLine():", string(l))
	log.Println(p.currentNode)
	log.Println(p.rootNode)
	log.Println("--End Nodes--")
	loop: for i := 0; i < len(l); i++ {
		log.Println("Line:", string(l), string(l[i]))
		switch l[i] {
			// Skip whitespace characters
			case ' ', '\t':
				break
			case '@':
				log.Println("Tag detected!:", string(l))
				// Is it a closing tag?
				if l[i+1] == '!' {
					name, _ := p.parseName(l[i+2:])
					log.Println("Closing node with name", name)
					if name != p.currentNode.Name {
						log.Fatalf("%d: Closing tag does not match currently open tag! Document is not well-formed! (got %s, expected %s)\n", p.lineNo, name, p.currentNode.Name)
					}
					p.currentNode = p.currentNode.Parent
					return
				}
				// Add this new node we are processing as a child
				log.Println(p.currentNode)
				log.Println(p.rootNode)
				newNode := NewSudNode(p.currentNode)
				p.currentNode.Children = append(p.currentNode.Children, newNode)
				p.currentNode = newNode
				
				p.inTag = true
				var n int
				// Get the current node name
				p.currentNode.Name, n = p.parseName(l[i+1:])
				p.parseAttributes(l[i+1+n:])
				// It was a single-line tag
				if !p.inTag {
					p.currentNode = p.currentNode.Parent
				}
				break loop
			default:
				p.currentNode.Content = append(p.currentNode.Content, l...)
				p.currentNode.Content = append(p.currentNode.Content, '\n')
				break loop
		}
	}
}

func (p *parseState) parseName(l []uint8) (name string, offset int) {
	defer fmt.Println("In parseName():", offset, name)
	nameStarted := false
	for i := 0; i < len(l); i++ {
		offset++
		switch l[i] {
			case ' ', '\t':
				l = l[i+1:]
				i = 0
				// This whitespace marks the end of the name.
				if nameStarted {
					return
				}
			default:
				nameStarted = true
				name += string(l[i])
		}
	}
	return
}

func (p *parseState) parseAttributes(l []uint8) {
	log.Println("In parseAttributes():", string(l))
	i := 0
	for i < len(l) {
		switch l[i] {
			case ' ', '\t':
				l = l[i+1:]
				i = 0
			default:
				i += p.parseAttribute(l[i:])
		}
		if !p.inTag {
			return
		}
	}
	return
}

func (p *parseState) parseAttribute(l []uint8) (i int) {
	fmt.Println("In parseAttribute():", string(l), len(l))
	val := ""
	key := ""
	parsingKey := true
	defer func() {
		if key != "" {
			fmt.Println("Found attribute:", key, "=", val)
			p.currentNode.Args[key] = val
		}
	}()
	for i = 0; i < len(l); i++ {
		//fmt.Println(string(l[i]))
		switch l[i] {
			case ' ', '\t':
				if p.inQuotes {
					if parsingKey {
						key += string(l[i])
					} else {
						val += string(l[i])
					}
				} else {
					// Next attribute
					return
				}
			case '"':
				// TODO: backslash escaping
				if p.inQuotes {
					p.inQuotes = false
				} else {
					p.inQuotes = true
				}
			case '=':
				parsingKey = false

			// End-of-tag markers
			case '\n':
				log.Fatalln("newline encountered in line-processed input. Wtf.")
				p.inTag = false
				return
			case '!':
				p.inTag = false
				return
			default:
				if parsingKey {
					key += string(l[i])
				} else {
					val += string(l[i])
				}
		}
	}
	return
}

func ParseReader(rd io.Reader) (root *Node, err os.Error) {
	bufr := bufio.NewReader(rd)
	ps := newParseState(bufr, os.Stdout)
	ps.rd = bufr
	ps.wr = os.Stdout
	for {
		l, isPrefix, err := bufr.ReadLine()
		if err != nil {
			if err == os.EOF {
				return ps.rootNode, nil
			}
			return nil, err
		}
		for isPrefix == true {
			var l2 []uint8
			l2, isPrefix, err = bufr.ReadLine()
			if err != nil {
				if err == os.EOF {
					return ps.rootNode, nil
				}
				return nil, err
			}
			l = append(l, l2...)
		}
		ps.parseLine(l)
		fmt.Println(string(l))
	}
	fmt.Println("Done parsing! Root node:", ps.rootNode)
	return ps.rootNode, nil
}