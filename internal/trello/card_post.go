package trello

import (
	"math"
	"net/url"
	"path"
	"strings"

	"github.com/ravbaker/fb2trello/internal/fb"
)

type cardPost fb.Post

func (p *cardPost) Name() string {
	name := strings.Join([]string{"FB Post", p.Message}, " - ")
	limit := int(math.Min(float64(len(name)), 512))
	return name[:limit]
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
		for _, attachment := range p.Attachments.Data {
			attachmentUrl := attachment.Url
			if len(attachmentUrl) > 0 {
				attachmentUrl = extractFbUrl(attachmentUrl, exists, links)
				if _, exist := exists[attachmentUrl]; !exist {
					exists[attachmentUrl] = true
					links = append(links, attachmentUrl)
				}
			}
		}
	}
	return
}

func extractFbUrl(attachmentUrl string, exists map[string]bool, links []string) string {
	parsedUrl, _ := url.Parse(attachmentUrl)
	extractedFBLink := parsedUrl.Query()["u"]
	if len(extractedFBLink) > 0 {
		attachmentUrl = extractedFBLink[0]
	}

	return attachmentUrl
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
