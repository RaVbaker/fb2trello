package fb

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var posts []Post

func ParseArchiveFolder() []Post {
	archive := sortedArchiveFiles()

	for _, entry := range archive {
		ProcessArchive(entry.Name())
	}

	log.Printf("Loaded %d posts from %d archives", len(posts), len(archive))
	return posts
}

func sortedArchiveFiles() []os.FileInfo {
	archive, err := ioutil.ReadDir(Folder)
	if err != nil {
		log.Fatalf("Cannot read `%s` folder", Folder)
	}
	sort.Slice(archive, func(i, j int) bool {
		name1, name2 := archive[i].Name(), archive[j].Name()
		if len(name1) == len(name2) {
			return strings.Compare(name1, name2) < 0
		}
		return len(name1) < len(name2)
	})
	return archive
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
