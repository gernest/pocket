package pocket

import (
	"time"

	"github.com/satori/go.uuid"
)

type Air struct {
	UUID  string         `json:"uuid"`
	Begin time.Time      `json:"begin"`
	End   time.Time      `json:"end"`
	Ads   map[string]*Ad `json:"ad"`
}

func NewAir() *Air {
	return &Air{
		UUID: uuid.NewV4().String(),
		Ads:  make(map[string]*Ad),
	}
}

func NewAIR(begin, end interface{}) (*Air, error) {
	a := &Air{}
	a.UUID = uuid.NewV4().String()
	switch t := begin.(type) {
	case string:
		n, err := time.Parse(time.RFC822, t)
		if err != nil {
			return nil, err
		}
		a.Begin = n
	case time.Time:
		a.Begin = t
	}
	switch t := end.(type) {
	case string:
		n, err := time.Parse(time.RFC822, t)
		if err != nil {
			return nil, err
		}
		a.End = n
	case time.Time:
		a.End = t
	}
	a.Ads = make(map[string]*Ad)
	return a, nil
}
func (a *Air) AddAd(adv *Ad) {
	a.Ads[adv.UUID] = adv
}
