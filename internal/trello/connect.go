package trello

import (
	"github.com/adlio/trello"
)

var client *trello.Client

func Connect(appKey, token string) {
	client = trello.NewClient(appKey, token)
}
