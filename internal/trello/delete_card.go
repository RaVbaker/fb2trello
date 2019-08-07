package trello

import (
	"fmt"
	"log"

	"github.com/adlio/trello"
)

func deleteCard(cardForDeletion *trello.Card) {
	err := client.Delete(fmt.Sprintf("cards/%s", cardForDeletion.ID), trello.Defaults(), &cardForDeletion)
	if err != nil {
		log.Printf("couldn't delete card %s", cardForDeletion.ID)
	}
}
