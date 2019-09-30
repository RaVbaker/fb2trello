package trello

import (
	"log"

	"github.com/ravbaker/fb2trello/internal/fb"

	"github.com/adlio/trello"
)

func cardFromPost(post *fb.Post) *trello.Card {
	postDetails := cardPost(*post)

	card := cardExists(&postDetails)
	if card != nil {
		return card
	}

	card = makeCard(&postDetails)

	// the post ID is used also for deduplication
	storeLink(card, &trello.Attachment{Name: "Post " + postDetails.Id, URL: postDetails.Link})

	for _, link := range postDetails.ExtractSharedLinks() {
		storeLink(card, &trello.Attachment{Name: "Link - " + link, URL: link})
	}

	if image := postDetails.ExtractImage(); len(image) > 0 {
		UploadAttachment(card, image)
	}
	log.Printf("Card: %s added successfully - %s", card.ID, card.Due)
	return card
}

func storeLink(card *trello.Card, link *trello.Attachment) {
	err := card.AddURLAttachment(link)
	if err != nil {
		deleteCard(card)
		log.Fatalf("Couldn't add card[%s] link %s attachment %s, details: %v", card.ID, link.Name, link.URL, err)
	}
}

func cardExists(postDetails *cardPost) *trello.Card {
	foundCards, err := client.SearchCards(postDetails.Id, trello.Defaults())
	if err != nil {
		log.Printf("Couldn't find cards, details: %v", err)
	} else if len(foundCards) > 0 {
		var card *trello.Card
		card, foundCards = foundCards[0], foundCards[1:]
		log.Printf("Found (%d), for query: `%s` - %s", len(foundCards)+1, postDetails.Id, card.Due)
		for _, cardForDeletion := range foundCards {
			deleteCard(cardForDeletion)
			log.Printf("Deleting extra card: %s", cardForDeletion.ID)
		}

		return card
	}
	return nil
}

func makeCard(postDetails *cardPost) *trello.Card {
	card := &trello.Card{
		Name:   postDetails.Name(),
		Desc:   postDetails.Desc(),
		IDList: list.ID,
	}
	err := list.AddCard(card, trello.Defaults())
	if err != nil || len(card.ID) == 0 {
		log.Fatalf("Couldn't add card[%s] `%s`, `%s`; details: %v", card.ID, card.Name, card.Desc, err)
	}
	// setting due and idLabels works only on Update, not on creation
	err = card.Update(trello.Arguments{"due": postDetails.CreatedTime, "dueComplete": "true", "idLabels": postDetails.Kind()})
	if err != nil {
		deleteCard(card)
		log.Fatalf("Card[%s] Due date set failed, details: %v", card.ID, err)
	}
	return card
}
