package fb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	url    = "https://graph.facebook.com/v3.3/%s/posts?fields=message%%2Cpicture%%2Cpermalink_url%%2Ccreated_time%%2Cfrom%%2Cfull_picture%%2Cmessage_tags%%2Cstatus_type%%2Cvia%%2Csharedposts%%7Blink%%7D%%2Cattachments&access_token=%s"
	Folder = "archive"
)

type JsonNextUrl struct {
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

var archiveCounter int = 1

func Archive(accessToken, pageName string) error {
	err := prepareArchiveFolder()
	fullUrl := fmt.Sprintf(url, pageName, accessToken)

	for len(fullUrl) != 0 {
		fmt.Printf("Fetching archive part: %d\n", archiveCounter)
		var body []byte
		body, err = storeRequest(fullUrl, fmt.Sprintf("%s/%d.json", Folder, archiveCounter))
		var nextLink JsonNextUrl
		err = json.Unmarshal(body, &nextLink)
		if err != nil {
			return err
		}

		archiveCounter++
		fullUrl = nextLink.Paging.Next
	}

	return nil
}

func prepareArchiveFolder() error {
	err := os.RemoveAll(Folder)
	if err != nil {
		return err
	}
	err = os.MkdirAll(Folder, os.ModePerm)
	return err
}

func storeRequest(url, destination string) ([]byte, error) {
	body, err := fetchUrl(url)
	if err != nil {
		return body, err
	}
	err = ioutil.WriteFile(destination, body, 0644)
	return body, err
}

func fetchUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	return body, err
}
