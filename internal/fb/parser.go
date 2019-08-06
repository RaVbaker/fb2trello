package fb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

var posts []Post

func ParseArchiveFolder() []Post {
	archive, err := ioutil.ReadDir(Folder)
	if err != nil {
		log.Fatalf("Cannot read `%s` folder", Folder)
	}
	for _, entry := range archive {
		ProcessArchive(entry.Name())
	}

	log.Printf("Loaded %d posts from %d archives", len(posts), len(archive))
	return posts
}

func ProcessArchive(fileName string) {
	var archive Result
	fullPath := filepath.Join(Folder, fileName)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Fatalf("Cannot read archive `%s` file", fullPath)
	}
	json.Unmarshal(content, &archive)

	posts = append(posts, archive.Data...)
}
