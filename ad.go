package pocket

import (
	"github.com/satori/go.uuid"
)

// Ad is the adevrtisement
type Ad struct {
	UUID string `json:"uuid"`
	Body string `json:"body"`
}

func NewAd(body string) *Ad {
	return &Ad{
		UUID: uuid.NewV4().String(),
		Body: body,
	}
}
