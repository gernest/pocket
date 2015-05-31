package pocket

import (
	"time"

	"github.com/robfig/cron"
)

// Ad is the adevrtisement
type Ad struct {
	UUID string `json:"uuid"`

	// Kind describes the type of the ad
	Kind int    `json:"kind"`
	Body string `json:"body"`
}

// Air is the air time of the ad
type Air struct {
	UUID  string    `json:"uuid"`
	Begin time.Time `json:"begin"`
	End   time.Time `json:"end"`
	Adv   Ad        `json:"ad"`
}

// Channel is the channel of ads
type Channel struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Airs      []*Air    `json:"airs"`
	Schedules *Schedule `json:"schedule"`
}

// BroadCast is where the channels of ads are
type BroadCast struct {
	UUID      string     `json:"uuid"`
	Channels  []*Channel `json:"channels"`
	Schedules *Schedule  `json:"schedule"`
}

// Schedule is the ad timings of broadcase.
type Schedule struct {

	// Pass are events which have already occured.
	Pass []*Event `jspo:"pass"`

	// Pending are events to be executed
	Pending []*Event `json:"pending"`

	events map[string]*Event `json:"-"`
	jobs   *cron.Cron        `json:"-"`
}

// Event is the occurance of events
type Event struct {
	UUID      string `json:"uuid"`
	Occurance string `json:"occurance"`
	EventID   string `json:"event_id"`
}
