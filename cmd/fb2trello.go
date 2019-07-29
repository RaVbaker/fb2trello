package main

import (
	"fmt"
	"os"
	
	"github.com/ravbaker/fb2trello/internal/fb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing argument: %s pageName", os.Args[0])
		os.Exit(1)
	}
	pageName := os.Args[1]
	err := fb.Archive(pageName)
	if err != nil {
		fmt.Errorf("Something wrong with FB archiving: %v", err)
		os.Exit(2)
	}
}
