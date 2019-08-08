package trello

import (
	"log"
	"strings"

	"github.com/adlio/trello"

	"github.com/ravbaker/fb2trello/internal/fb"
)

var board *trello.Board
var list *trello.List
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
	board = getBoard(boardName)
	list = getList(board, listName)
	prepareLabels()
	log.Printf("Board[%s] `%s` ready: %s", board.ID, board.Name, board.URL)
}
