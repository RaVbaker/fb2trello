package trello

import (
	"log"

	"github.com/adlio/trello"
)

func SetupBoard(boardName string, lists []string) {
	existingBoard := getBoard(boardName)
	if existingBoard == nil {
		makeBoard(existingBoard, boardName)
	}
	makeLists(existingBoard, &lists)
}

func makeLists(board *trello.Board, lists *[]string) {
	var listToName = make(map[string]*trello.List)
	for _, list := range getLists(board) {
		listToName[list.Name] = list
	}

	for _, listName := range *lists {
		if _, exists := listToName[listName]; !exists {
			makeList(board, listName)
		}
	}
}

func makeList(board *trello.Board, listName string) {
	_, err := board.CreateList(listName, trello.Arguments{"pos": "bottom"})
	if err != nil {
		log.Fatalf("Failed to create list[%s], details: %v", listName, err)
	}
	log.Printf("Creating new list %s", listName)
}

func makeBoard(board *trello.Board, boardName string) {
	board.Name = boardName
	err := client.CreateBoard(board, trello.Arguments{"powerUps": "calendar"})
	if err != nil {
		log.Fatalf("Could not create board %s, details: %v", boardName, err)
	}
	log.Printf("Creating new board %s", boardName)
}

func getBoard(boardName string) *trello.Board {
	boards, err := client.GetMyBoards(trello.Arguments{"lists": "open"})
	if err != nil {
		log.Fatalf("Cannot retrieve Trello Boards, details: %v", err)
	}
	for _, board := range boards {
		if board.Name == boardName {
			return board
		}
	}
	return nil
}

func getList(board *trello.Board, listName string) *trello.List {
	lists := getLists(board)
	
	for _, list := range lists {
		if list.Name == listName {
			return list
		}
	}
	return nil
}

func getLists(board *trello.Board) []*trello.List {
	log.Printf("Fetch lists from board %s", board.ID)
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatalf("Cannot retrieve Trello Board[%s] lists, details: %v", board.ID, err)
	}
	return lists
}
