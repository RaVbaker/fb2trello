package trello

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/ravbaker/fb2trello/internal/fb"

	"github.com/adlio/trello"
)

type cardPost fb.Post

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

func deleteCard(cardForDeletion *trello.Card) {
	err := client.Delete(fmt.Sprintf("cards/%s", cardForDeletion.ID), trello.Defaults(), &cardForDeletion)
	if err != nil {
		log.Printf("couldn't delete card %s", cardForDeletion.ID)
	}
}

func makeCard(postDetails *cardPost) *trello.Card {
	card := &trello.Card{
		Name:   postDetails.Name(),
		Desc:   postDetails.Desc(),
		IDList: list.ID,
	}
	err := list.AddCard(card, trello.Defaults())
	if err != nil || len(card.ID) == 0 {
		log.Fatalf("Couldn't add card[%s] `%s`, details: %v", card.ID, card.Name, err)
	}
	// setting due and idLabels works only on Update, not on creation
	err = card.Update(trello.Arguments{"due": postDetails.CreatedTime, "dueComplete": "true", "idLabels": postDetails.Kind()})
	if err != nil {
		deleteCard(card)
		log.Fatalf("Card[%s] Due date set failed, details: %v", card.ID, err)
	}
	return card
}

func (p *cardPost) Name() string {
	return strings.Join([]string{"FB Post", p.Message}, " - ")
}

func (p *cardPost) Desc() string {
	var attachmentContent string
	if len(p.Attachments.Data) > 0 {
		attachment := p.Attachments.Data[0]
		attachmentContent = attachment.Title + "\n" + attachment.Description
	}
	links := p.ExtractSharedLinks()
	content := []string{
		p.Message,
		"",
		attachmentContent,
		"",
		"Time:",
		p.CreatedTime,
		"",
		"Post:",
		p.Link,
		"",
		"Links:",
	}
	content = append(content, links...)
	return strings.Join(content, "\n")
}

func (p *cardPost) ExtractSharedLinks() (links []string) {
	exists := make(map[string]bool)
	for _, share := range p.SharedPosts.Data {
		if _, exist := exists[share.Link]; !exist {
			exists[share.Link] = true
			links = append(links, share.Link)
		}
	}
	if len(p.Attachments.Data) > 0 {
		if attachmentUrl := p.Attachments.Data[0].Url; len(attachmentUrl) > 0 {
			parsedUrl, _ := url.Parse(attachmentUrl)
			extractedFBLink := parsedUrl.Query()["u"]
			if len(extractedFBLink) > 0 {
				attachmentUrl = extractedFBLink[0]
			}
			if _, exist := exists[attachmentUrl]; !exist {
				exists[attachmentUrl] = true
				links = append(links, attachmentUrl)
			}
		}
	}
	return
}

func (p *cardPost) Kind() string {
	if p.StatusType == "added_photos" {
		return labels["Photo"]
	}

	if len(p.Attachments.Data) > 0 {
		attachment := p.Attachments.Data[0]
		if attachment.Type == "photo" {
			return labels["Photo"]
		} else if attachment.Type == "video_inline" || strings.Contains(attachment.Url, "youtube") {
			return labels["Video"]
		}
		return labels["Link"]
	}

	return labels["Post"]
}

func (p *cardPost) ExtractImage() string {
	imageUrl := p.Picture
	if len(p.Attachments.Data) > 0 {
		if image := p.Attachments.Data[0].Media.Src; len(image) > 0 {
			imageUrl = image
		}
	}

	parsedUrl, _ := url.Parse(imageUrl)
	if path.Ext(strings.Split(parsedUrl.Path, "?")[0]) == ".php" {
		imageUrl = parsedUrl.Query()["url"][0]
	}
	return imageUrl
}
