package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ravbaker/fb2trello/internal/fb"
	"github.com/ravbaker/fb2trello/internal/trello"
)

func main() {
	var pageName, boardName, listNames, untilDate string
	var setup bool
	flag.StringVar(&pageName, "page", "", "Facebook pageName/ID  which should get archived")
	flag.StringVar(&boardName, "board", "", "Trello board name to which posts should be archived")
	flag.StringVar(&listNames, "lists", "Calendar,Ideas,Planned,Published", "Trello list names, last after comma is the one for archive, default/e.g. `Calendar,Ideas,Planned,Published`")
	flag.StringVar(&untilDate, "until", "", "Archive until date - oldest post publication date, e.g. 2019-07-30")
	flag.BoolVar(&setup, "setup", false, "If specified it will create whole structure in trello for board and lists")
	flag.Parse()

	if len(pageName) > 0 {
		accessToken := getEnvVar("FB_ACCESS_TOKEN")

		err := fb.Archive(accessToken, pageName, untilDate)
		if err != nil {
			fmt.Errorf("Something wrong with FB archiving: %v", err)
			os.Exit(2)
		}
	}

	if len(boardName) > 0 {
		apiKey := getEnvVar("TRELLO_API_KEY")
		token := getEnvVar("TRELLO_TOKEN")
		lists := strings.Split(listNames, ",")

		trello.Connect(apiKey, token)
		if setup {
			trello.SetupBoard(boardName, lists)
		}
		posts := fb.ParseArchiveFolder()
		lastList := lists[len(lists)-1]
		trello.StoreCards(boardName, lastList, untilDate, posts)
	}
}

func getEnvVar(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("%s ENV variable not found. Please provide!", name)
	}
	return value
}
