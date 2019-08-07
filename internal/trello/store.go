package trello

import (
	"fmt"
	"log"
	"strings"

	"github.com/adlio/trello"

	"github.com/ravbaker/fb2trello/internal/fb"
)

var board trello.Board
var list trello.List
var labels = map[string]string{
	"Post":  "",
	"Photo": "",
	"Video": "",
	"Link":  "",
}

func StoreCards(boardName, listName, untilDate string, posts []fb.Post) {
	loadBoardDetails(boardName, listName)

	var counter uint
	for _, post := range posts {
		if len(untilDate) > 0 && strings.Compare(post.CreatedTime, untilDate) < 0 {
			break
		}
		cardFromPost(&post)
		counter++
	}
	log.Printf("Processed %d posts", counter)
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
	prepareLabels()
	log.Printf("Board[%s] `%s` ready: %s", board.ID, board.Name, board.URL)
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

func prepareLabels() {
	loadedLabels, err := board.GetLabels(trello.Defaults())
	if err != nil {
		log.Fatalf("Cannot retrieve Trello Board loadedLabels, details: %v", err)
	}

	availableLabels := discoverExistingLabels(loadedLabels)
	addMissingLabels(availableLabels, loadedLabels)
}

func discoverExistingLabels(loadedLabels []*trello.Label) []int {
	var availableLabels []int
	for i, boardLabel := range loadedLabels {
		if _, exists := labels[boardLabel.Name]; exists {
			log.Printf("label[%s] found with name: `%v`", boardLabel.ID, boardLabel.Name)
			labels[boardLabel.Name] = boardLabel.ID
		} else if len(boardLabel.Name) == 0 {
			availableLabels = append(availableLabels, i)
		}
	}
	return availableLabels
}

func addMissingLabels(availableLabels []int, loadedLabels []*trello.Label) {
	var label *trello.Label
	for name, id := range labels {
		// skip labels with IDs assigned
		if len(id) > 0 {
			continue
		}

		if len(availableLabels) > 0 {
			label, availableLabels = loadedLabels[availableLabels[0]], availableLabels[1:]
			updateLabelName(name, label)
			labels[label.Name] = label.ID
			log.Printf("label[%s] %s updated: %s", label.ID, label.Name)
		} else {
			log.Fatalf("No available labels (with empty labels). Make some new and run script again")
		}
	}
}

func updateLabelName(name string, label *trello.Label) {
	args := trello.Defaults()
	args["value"] = name
	err := client.Put(fmt.Sprintf("labels/%s/name", label.ID), args, &label)
	if err != nil {
		log.Fatalf("Cannot update label[%s], name: %s, details: %v", label.ID, name, err)
	}
}
