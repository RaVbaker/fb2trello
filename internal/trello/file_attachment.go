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
	"github.com/pkg/errors"
)

func UploadAttachment(card *trello.Card, url string) *trello.Attachment {
	client.Throttle()

	resp, err := http.Get(url)
	if err != nil {
		deleteCard(card)
		log.Fatalf("card[%s] cannot read attachment %s, details: %v", card.ID, url, err)
	}
	defer resp.Body.Close()

	postData := &bytes.Buffer{}
	writer, err := constructMutlipartBody(url, postData, card, &resp.Body)
	if err != nil {
		deleteCard(card)
		log.Fatalf("card[%s] %s", card.ID, err.Error())
	}
	return createAttachmentUpload(card, postData, writer.FormDataContentType(), url)
}

func constructMutlipartBody(url string, postData *bytes.Buffer, card *trello.Card, respBody *io.ReadCloser) (*multipart.Writer, error) {
	basePath := path.Base(strings.Split(url, "?")[0])
	writer := multipart.NewWriter(postData)
	part, err := writer.CreateFormFile("file", basePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot write attachment %s - %s", url, basePath)
	}
	_, err = io.Copy(part, *respBody)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot copy attachment %s", basePath)
	}
	ext := imageExtension(basePath)
	err = writer.WriteField("mimeType", mime.TypeByExtension(ext))
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot write mimetype for %s [%s]", ext, url)
	}
	err = writer.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "Closing writer for %s", url)
	}
	return writer, nil
}

func imageExtension(basePath string) string {
	ext := path.Ext(basePath)
	if len(ext) < 2 {
		ext = ".jpg" // default extension
	}
	return ext
}

func createAttachmentUpload(card *trello.Card, postData *bytes.Buffer, contentType, url string) *trello.Attachment {
	params := urlPath.Values{"key": {client.Key}, "token": {client.Token}}
	callPath := strings.Join([]string{client.BaseURL, "cards", card.ID, "attachments"}, "/") + "?" + params.Encode()
	req, err := http.NewRequest("POST", callPath, postData)
	if err != nil {
		deleteCard(card)
		log.Fatalf("card[%s] Cannot make POST to upload attachment:%s, details %v", card.ID, callPath, err)
	}
	req.Header.Set("Content-Type", contentType)
	attachmentResp, err := client.Client.Do(req)
	if err != nil {
		deleteCard(card)
		log.Fatalf("card[%s]  HTTP request failure on %s, details %v", card.ID, url, err)
	}
	defer attachmentResp.Body.Close()
	return parseAttachmentResponse(attachmentResp, card, url)
}

func parseAttachmentResponse(attachmentResp *http.Response, card *trello.Card, url string) *trello.Attachment {
	b, err := ioutil.ReadAll(attachmentResp.Body)
	if err != nil {
		deleteCard(card)
		log.Fatalf("card[%s] HTTP Read error on response for %s, details %v", card.ID, url, err)
	}
	newAttachment := &trello.Attachment{}
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	err = decoder.Decode(newAttachment)
	if err != nil {
		deleteCard(card)
		log.Fatalf("card[%s] JSON decode failed on %s, details: %v:\n%s", card.ID, url, err, string(b))
	}
	return newAttachment
}
