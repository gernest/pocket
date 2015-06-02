package pocket

import (
	"testing"
	"time"
)

var (
	now         = time.Now()
	channelName = "101"
	dataSet     = []struct {
		uuid  string
		msg   string
		begin time.Time
		end   time.Time
	}{
		{
			"c1eeec26-901d-4c1a-99d0-8413fa6b8484", "hello worlf",
			now, now.Add(time.Second),
		},
		{
			"20867ddd-325b-428a-8d49-5f338a430b04", "mambo jambo",
			now.Add(time.Second), now.Add(2 * time.Second),
		},
		{"8b6a81f5-1384-4dbe-ace4-b01d6d318324", "habari gani",
			now.Add(2 * time.Second), now.Add(3 * time.Second),
		},
	}
)

func TestChannel_AddAirTime(t *testing.T) {
	ch := NewChannel("101")
	for _, v := range dataSet {
		air, err := NewAIR(v.begin, v.end)
		if err != nil {
			t.Errorf("creating air: %v", err)
		}
		air.UUID = v.uuid
		for i := range iter(3) {
			ad := NewAd(dataSet[i].msg)
			air.AddAd(ad)
		}
		ch.AddAirTime(air)
	}
}

func TestGetChannel(t *testing.T) {
	ch, err := GetChannel(channelName)
	if err != nil {
		t.Errorf("getting channel :%v", err)
	}
	if ch.Name != channelName {
		t.Errorf("expected %s got %s", channelName, ch.Name)
	}
	ch, err = GetChannel("")
	if err == nil {
		t.Error("expected an error got nil")
	}
	if ch != nil {
		t.Errorf("expected nil got %v", ch)
	}
}

func TestChannel_CurrentAirTime(t *testing.T) {
	defer db.DeleteDatabase()
	ch, err := GetChannel(channelName)
	if err != nil {
		t.Errorf("getting channel :%v", err)
	}
	if ch.Name != channelName {
		t.Errorf("expected %s got %s", channelName, ch.Name)
	}
	curr := ch.CurrentAirTime()
	if curr.UUID != dataSet[0].uuid {
		t.Errorf("expected %v go %v", dataSet[0].uuid, curr.UUID)
	}

	time.Sleep(2 * time.Second)
	curr = ch.CurrentAirTime()
	if curr.UUID != dataSet[1].uuid {
		t.Errorf("expected %s go %s", dataSet[1].uuid, curr.UUID)
	}

	time.Sleep(time.Second)
	curr = ch.CurrentAirTime()
	if curr.UUID != dataSet[2].uuid {
		t.Errorf("expected %s go %s", dataSet[2].uuid, curr.UUID)
	}
}

func iter(n int) []struct{} {
	return make([]struct{}, n)
}
