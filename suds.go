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
	root, err := Parse(f)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(root)
	log.Println("-------Generating HTML Output-------")
	err = GenerateHTML(root, os.Stdout)
	if err != nil {
		log.Println(err)
	}
}