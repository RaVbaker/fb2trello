package trello

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	urlPath "net/url"
	"path"
	"strings"

	"github.com/adlio/trello"
)

func UploadAttachment(card *trello.Card, url string) *trello.Attachment {
	client.Throttle()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("cannot read attachment %s, details: %v", url, err)
	}
	defer resp.Body.Close()

	basePath := path.Base(strings.Split(url, "?")[0])

	post := &bytes.Buffer{}
	writer := multipart.NewWriter(post)
	part, err := writer.CreateFormFile("file", basePath)
	if err != nil {
		log.Fatalf("cannot write attachment %s - %s, details: %v", url, basePath, err)
	}

	_, err = io.Copy(part, resp.Body)
	if err != nil {
		log.Fatalf("cannot copy attachment %s, details: %v", basePath, err)
	}

	ext := path.Ext(basePath)
	if len(ext) < 2 {
		ext = ".jpg" // default extension
	}
	err = writer.WriteField("mimeType", mime.TypeByExtension(ext))
	if err != nil {
		log.Fatalf("Cannot write mimetype for %s [%s], details %v", url, ext, err)
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("Closing writer for %s, details %v", url, err)
	}

	params := urlPath.Values{"key": {client.Key}, "token": {client.Token}}
	callPath := strings.Join([]string{client.BaseURL, "cards", card.ID, "attachments"}, "/") + "?" + params.Encode()
	req, err := http.NewRequest("POST", callPath, post)
	if err != nil {
		log.Fatalf("Cannot make POST to upload attachment:%s, details %v", callPath, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	attachmentResp, err := client.Client.Do(req)
	if err != nil {
		log.Fatalf("HTTP request failure on %s, details %v", url, err)
	}
	defer attachmentResp.Body.Close()

	b, err := ioutil.ReadAll(attachmentResp.Body)
	if err != nil {
		log.Fatalf("HTTP Read error on response for %s, details %v", url, err)
	}

	newAttachment := &trello.Attachment{}
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	err = decoder.Decode(newAttachment)
	if err != nil {
		log.Fatalf("JSON decode failed on %s, details: %v:\n%s", url, err, string(b))
	}
	return newAttachment
}
