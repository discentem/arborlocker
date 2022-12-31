package main

import (
	"fmt"
	"log"

	"github.com/discentem/arborlocker/github/pullrequest"
)

func main() {
	for i := 348; i < 349; i++ {
		pr, err := pullrequest.Query(nil, "facebook", "sapling", i)
		if err != nil {
			log.Print(err)
		}
		b := string(pr.Repository.PullRequest.BodyHTML)
		fmt.Print(b)
		fmt.Println(pullrequest.LinesFromHTMLDescription(b))
	}

}
