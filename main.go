package main

import (
	"log"

	"github.com/discentem/arborlocker/github/pullrequest"
)

func main() {
	pr, err := pullrequest.Query(nil, "facebook", "sapling", 348)
	if err != nil {
		log.Print(err)
	}

	// log.Println("server started")
	// http.HandleFunc("/webhook", pullrequest.RunWebhook)
	// log.Fatal(http.ListenAndServe(":3000", nil))

}
