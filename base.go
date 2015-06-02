package pocket

// Air is the air time of the ad

// BroadCast is where the channels of ads are
type BroadCast struct {
	UUID     string     `json:"uuid"`
	Channels []*Channel `json:"channels"`
}
