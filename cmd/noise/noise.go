package main

import (
	"fmt"
	"log"

	"github.com/pkar/noise"
)

func main() {
	n, err := noise.New()
	if err != nil {
		log.Fatal(err)
	}
	err = n.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
