package pocket

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/gernest/nutz"
	"github.com/satori/go.uuid"
)

const (
	channelBucket   = "channels"
	broadcastBucket = "broadcasts"
)

var (
	errScheduleNotFound = errors.New("pocket: schedule not found")
	errAirtimeExists    = errors.New("pocket: airtime already exists")
)

// Ad is the adevrtisement
type Ad struct {
	UUID string `json:"uuid"`
	Body string `json:"body"`
}

// Air is the air time, its a time slot or in which particular ads reside.
type Air struct {
	UUID  string         `json:"uuid"`
	Begin time.Time      `json:"begin"`
	End   time.Time      `json:"end"`
	Ads   map[string]*Ad `json:"ad"`
}

// BroadCast is where the channels of ads are
type BroadCast struct {
	UUID     string   `json:"uuid"`
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	db       nutz.Storage
}

// Channel is the channel of ads
type Channel struct {
	UUID      string          `json:"uuid"`
	Name      string          `json:"name"`
	AirTime   map[string]*Air `json:"air_time"`
	BroadCast string          `json:"broadcast"`
	timeTable *Scheduler      `jsson:"-"`
	db        nutz.Storage
	mutex     sync.RWMutex
}

type schedule struct {
	start    time.Time
	duration time.Duration
	event    interface{}
}

// Scheduler a tme based schedule of events.
type Scheduler struct {
	table schedules
	mux   sync.RWMutex
}

type schedules []*schedule

// NewAd creates a new adverisement, it assigns a new UUID to the ad.
func NewAd(body string) *Ad {
	return &Ad{
		UUID: uuid.NewV4().String(),
		Body: body,
	}
}

// NewAir creates a new Air time and assigns to it a new UUID
func NewAir(begin, end interface{}) (*Air, error) {
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

// AddAd adds an Ad to Air.
func (a *Air) AddAd(adv *Ad) {
	a.Ads[adv.UUID] = adv
}

// NewChannel creates a new channel, it assigns a new UUID.
func NewChannel(name, cast string, store nutz.Storage) *Channel {
	return &Channel{
		UUID:      uuid.NewV4().String(),
		Name:      name,
		AirTime:   make(map[string]*Air),
		BroadCast: cast,
		timeTable: NewScheduler(),
		db:        store,
	}
}

// AddAirTime adds a new airtime to the channel
func (c *Channel) AddAirTime(a *Air) error {
	if c.Exists(a) {
		return errAirtimeExists
	}
	c.mutex.RLock()
	c.AirTime[a.UUID] = a
	c.timeTable.Add(&schedule{
		start:    a.Begin,
		duration: a.End.Sub(a.Begin),
		event:    a,
	})
	c.mutex.RUnlock()
	err := c.Save()
	if err != nil {
		panic(err)
	}
	return nil
}

// Exists checks whether the Air time exist.
func (c *Channel) Exists(a *Air) bool {
	for _, v := range c.AirTime {
		if v.Begin.Equal(a.Begin) {
			return true
		}
	}
	return false
}

// CurrentAirTime retrieves what is currently on air.
func (c *Channel) CurrentAirTime() (*Air, error) {
	s, err := c.timeTable.onAir()
	if err != nil {
		return nil, err
	}
	return s.event.(*Air), nil
}

// Save persist the channel to bolt database.
func (c *Channel) Save() error {
	d, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return c.db.Create(channelBucket, c.Name, d, c.BroadCast).Error
}

// loads the airtimes to the scheduler
func (c *Channel) load() {
	for _, a := range c.AirTime {
		c.timeTable.Add(&schedule{
			start:    a.Begin,
			duration: a.End.Sub(a.Begin),
			event:    a,
		})
	}
}

// GetChannel retrieves a channel from the database and loads the aitimes.
func GetChannel(name, cast string, store nutz.Storage) (*Channel, error) {
	c := store.Get(channelBucket, name, cast)
	if c.Error != nil {
		return nil, c.Error
	}
	ch := NewChannel(name, cast, store)
	err := json.Unmarshal(c.Data, ch)
	if err != nil {
		return nil, err
	}
	ch.load()
	return ch, nil
}

// NewBroadcast creates a new broadcast
func NewBroadcast(name string, store nutz.Storage) *BroadCast {
	return &BroadCast{
		UUID: uuid.NewV4().String(),
		Name: name,
		db:   store,
	}
}

// GetBroadCast retrieves a broadcast
func GetBroadCast(name string, store nutz.Storage) (*BroadCast, error) {
	b := store.Get(broadcastBucket, name)
	if b.Error != nil {
		return nil, b.Error
	}
	cast := &BroadCast{}
	err := json.Unmarshal(b.Data, cast)
	if err != nil {
		return nil, err
	}
	return cast, nil
}

// HasChannel checks where the broadcast has a channel with the ginen name.
func (b *BroadCast) HasChannel(s string) bool {
	for _, v := range b.Channels {
		if v == s {
			return true
		}
	}
	return false
}

// AddChannel ads a channel to the broadcast
func (b *BroadCast) AddChannel(ch *Channel) {
	b.Channels = append(b.Channels, ch.Name)
	ch.BroadCast = b.Name
	b.Save()
	ch.Save()
}

// Save saves a broadcast
func (b *BroadCast) Save() error {
	d, err := json.Marshal(b)
	if err != nil {
		return err
	}
	c := b.db.Create(broadcastBucket, b.Name, d)
	return c.Error
}

// A sorted list of schedules implementation
func (s schedules) Len() int           { return len(s) }
func (s schedules) Less(i, j int) bool { return s[i].start.Before(s[j].start) }
func (s schedules) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		table: make(schedules, 0),
	}
}

// Add adds a schedule to a scheduler
func (s *Scheduler) Add(sh *schedule) {
	//s.table = append(s.table, sh)
	s.mux.RLock()
	s.table = append(s.table, sh)
	s.mux.RUnlock()
}

// returns a current schedule.
func (s *Scheduler) onAir() (*schedule, error) {
	sort.Sort(s.table)
	return s.table.Now()
}

// Now returns what schedule to show right now.
func (s schedules) Now() (*schedule, error) {
	now := time.Now()
	i := sort.Search(len(s), func(n int) bool {
		shd := s[n]
		if shd.start.Equal(now) {
			return true
		}
		if shd.start.Before(now) && now.Sub(shd.start) < shd.duration {
			return true
		}
		return false
	})
	if i == len(s) {
		return nil, errScheduleNotFound
	}
	rst := s[i]
	if rst.start.After(now) {
		return nil, errScheduleNotFound
	}
	return rst, nil
}
