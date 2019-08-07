package trello

import (
	"net/url"
	"path"
	"strings"

	"github.com/ravbaker/fb2trello/internal/fb"
)

type cardPost fb.Post

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
