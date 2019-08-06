package trello

import (
	"log"

	"github.com/adlio/trello"

	"github.com/ravbaker/fb2trello/internal/fb"
)

var board trello.Board
var list trello.List
var labels map[string]string

func StoreCards(boardName, listName string, posts []fb.Post) {
	loadBoardDetails(boardName, listName)

	log.Printf("My board %v", board)

	for _, post := range posts {
		cardFromPost(&post)
	}
}

func loadBoardDetails(boardName, listName string) {
	boards, err := client.GetMyBoards(trello.Defaults())
	if err != nil {
		log.Fatalf("Cannot retrieve Trello Boards, details: %v", err)
	}
	for _, myBoard := range boards {
		if myBoard.Name == boardName {
			board = *myBoard
			break
		}
	}

	loadList(listName)
	loadLabels()
}

func loadList(listName string) {
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatalf("Cannot retrieve Trello Board lists, details: %v", err)
	}
	for _, boardList := range lists {
		if boardList.Name == listName {
			list = *boardList
			break
		}
	}
}

func loadLabels() {
	labels = make(map[string]string)
	loadedLabels, err := board.GetLabels(trello.Defaults())
	if err != nil {
		log.Fatalf("Cannot retrieve Trello Board loadedLabels, details: %v", err)
	}
	for _, boardLabel := range loadedLabels {
		if len(boardLabel.Name) > 0 {
			log.Printf("name `%v`, id: %s", boardLabel.Name, boardLabel.ID)
			labels[boardLabel.Name] = boardLabel.ID
		}
	}
}
