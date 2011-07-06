package main

import (
	"os"
	"log"
)

func main() {
	f, err := os.Open("example1.suds")
	if err != nil {
		log.Fatalln(err)
	}
	root, err := ParseReader(f)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(root)
}