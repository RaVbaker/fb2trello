package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ravbaker/fb2trello/internal/fb"
	"github.com/ravbaker/fb2trello/internal/trello"
)

func main() {
	var pageName, boardName, listName string
	flag.StringVar(&pageName, "page", "", "Facebook pageName/ID  which should get archived")
	flag.StringVar(&boardName, "board", "", "Trello board name to which posts should be archived")
	flag.StringVar(&listName, "list", "", "Trello list name to which posts should be archived")
	flag.Parse()

	if len(pageName) > 0 {
		accessToken := getEnvVar("FB_ACCESS_TOKEN")

		err := fb.Archive(accessToken, pageName)
		if err != nil {
			fmt.Errorf("Something wrong with FB archiving: %v", err)
			os.Exit(2)
		}
	}

	if len(boardName) > 0 {
		apiKey := getEnvVar("TRELLO_API_KEY")
		token := getEnvVar("TRELLO_TOKEN")

		posts := fb.ParseArchiveFolder()
		trello.Connect(apiKey, token)
		trello.StoreCards(boardName, listName, posts[:50])
	}
}

func getEnvVar(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("%s ENV variable not found. Please provide!", name)
	}
	return value
}
