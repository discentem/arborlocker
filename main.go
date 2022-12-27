package main

import (
	"fmt"
	"log"

	"github.com/discentem/arborlocker/github/pullrequest"
)

func main() {
	pr, err := pullrequest.Query(nil, 348)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pr)
}
