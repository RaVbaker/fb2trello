package trello

import (
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

	foundCards, err := client.SearchCards(postDetails.Id, trello.Defaults())
	if err != nil {
		log.Printf("Couldn't find cards, details: %v", err)
	} else if len(foundCards) > 0 {
		log.Printf("Found cards(%d), query: `%s`: %v", len(foundCards), postDetails.Id, foundCards)
		return foundCards[0]
	}

	card := &trello.Card{
		Name:   postDetails.Name(),
		Desc:   postDetails.Desc(),
		IDList: list.ID,
	}
	err = list.AddCard(card, trello.Defaults())
	if err != nil || len(card.ID) == 0 {
		log.Fatalf("Couldn't add card[%s] `%s`, details: %v", card.ID, card.Name, err)
	}

	// setting due and idLabels works only on Update, not on creation
	err = card.Update(trello.Arguments{"due": postDetails.CreatedTime, "dueComplete": "true", "idLabels": postDetails.Kind()})
	if err != nil {
		log.Fatalf("Card[%s] Due date set failed, details: %v", card.ID, err)
	}

	// the post ID is used also for deduplication
	err = card.AddURLAttachment(&trello.Attachment{Name: "Post " + postDetails.Id, URL: postDetails.Link})
	if err != nil {
		log.Fatalf("Couldn't add card `%s` attachment %s, details: %v", card.Name, postDetails.Link, err)
	}

	for pos, link := range postDetails.ExtractSharedLinks() {
		err = card.AddURLAttachment(&trello.Attachment{Name: "Link - " + link, URL: link})
		if err != nil {
			log.Fatalf("Couldn't add card attachment[%d] %s, details: %v", pos, link, err)
		}
	}

	if image := postDetails.ExtractImage(); len(image) > 0 {
		UploadAttachment(card, image)
	}
	log.Printf("Card: %s added successfully", card.ID)
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
		} else if attachment.Type == "video_inline" {
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
