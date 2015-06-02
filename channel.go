package pocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gernest/nutz"
	"github.com/ryszard/goskiplist/skiplist"
	"github.com/satori/go.uuid"
)

const (
	channelBucket   = "channels"
	airTimeInterval = 2 * time.Second
)

var db = nutz.NewStorage("pocket.db", 0600, nil)

// Channel is the channel of ads
type Channel struct {
	UUID    string             `json:"uuid"`
	Name    string             `json:"name"`
	AirTime map[string]*Air    `json:"air_time"`
	set     *skiplist.SkipList `jsson:"-"`
	mutex   sync.RWMutex
}

func NewChannel(name string, id ...string) *Channel {
	return &Channel{
		UUID:    uuid.NewV4().String(),
		Name:    name,
		AirTime: make(map[string]*Air),
		set: skiplist.NewCustomMap(func(l, r interface{}) bool {
			return l.(int64) < r.(int64)
		}),
	}
}

func (c *Channel) AddAirTime(a *Air) {
	c.mutex.Lock()
	c.AirTime[a.UUID] = a
	c.set.Set(a.End.Unix(), a)
	c.mutex.Unlock()
	go c.Save()
}

func (c *Channel) CurrentAirTime() *Air {
	t := time.Now().Add(airTimeInterval)
	r := c.set.Range(time.Now().Unix(), t.Unix())

	var v []*Air
	for r.Next() {
		v = append(v, r.Value().(*Air))
	}
	if v != nil {
		return v[0]
	}
	return nil
}
func (c *Channel) Save() error {
	d, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return db.Create(channelBucket, c.Name, d).Error
}

func (c *Channel) load() {
	if c.set == nil {
		c.set = skiplist.NewCustomMap(func(l, r interface{}) bool {
			return l.(int64) < r.(int64)
		})
	}
	for _, v := range c.AirTime {
		c.set.Set(v.End.Unix(), v)
	}
}

func GetChannel(name string) (*Channel, error) {
	c := db.Get(channelBucket, name)
	if c.Error != nil {
		return nil, c.Error
	}
	ch := &Channel{}
	err := json.Unmarshal(c.Data, ch)
	if err != nil {
		return nil, err
	}
	ch.load()
	return ch, nil
}
