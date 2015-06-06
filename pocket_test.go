package pocket

import (
	"testing"
	"time"

	"github.com/gernest/nutz"
)

var (
	testDb = nutz.NewStorage("pocket.db", 0600, nil)
	cast   = "home"

	duration = time.Millisecond
	now      = time.Now()
	data     = []struct {
		evt   string
		start time.Time
	}{
		{"hello world", now.Add(-duration)},
		{"hello africa", now},
		{"hello Tanzania", now.Add(time.Second)},
	}
	channelName = "101"
	dataSet     = []struct {
		uuid  string
		msg   string
		begin time.Time
		end   time.Time
	}{
		{
			"c1eeec26-901d-4c1a-99d0-8413fa6b8484", "hello worlf",
			now.Add(-duration), now,
		},
		{
			"20867ddd-325b-428a-8d49-5f338a430b04", "mambo jambo",
			now, now.Add(duration),
		},
		{"8b6a81f5-1384-4dbe-ace4-b01d6d318324", "habari gani",
			now.Add(duration), now.Add(3 * time.Millisecond),
		},
	}
)

func TestNewAir(t *testing.T) {
	begin := time.Now()
	end := begin.Add(time.Minute)

	a1, err := NewAir(begin, end)
	if err != nil {
		t.Errorf("creating new air: %v", err)
	}
	if !a1.Begin.Equal(begin) {
		t.Errorf("expected %s got %s", begin, a1.Begin)
	}
	a2, err := NewAir(begin.String(), end.String())
	if err == nil {
		t.Errorf("expected an error %v", err)
	}
	if a2 != nil {
		t.Errorf("expected nil got %v", a2)
	}
	a3, err := NewAir(begin.Format(time.RFC822), end.Format(time.RFC822))
	if err != nil {
		t.Errorf("creating new air: %v", err)
	}
	if a3.Begin.Format(time.RFC822) != begin.Format(time.RFC822) {
		t.Errorf("expected %s got %s", begin.Format(time.RFC822), a3.Begin.Format(time.RFC822))
	}
	a4, err := NewAir(begin, end.String())
	if err == nil {
		t.Errorf("expected an error %v", err)
	}
	if a4 != nil {
		t.Errorf("expected nil got %v", a2)
	}
}

func TestScheduler(t *testing.T) {
	sh := NewScheduler()
	for _, v := range data {
		shd := &schedule{
			start:    v.start,
			duration: duration,
			event:    v.evt,
		}
		sh.Add(shd)
	}
	rst, err := sh.OnAir()
	if err != nil {
		t.Errorf("getting current schedule: %v", err)
	}
	if rst.event != data[1].evt {
		t.Errorf("expected %s got %s", data[1].evt, rst.event)
	}
}

func TestChannel_AddAirTime(t *testing.T) {
	ch := NewChannel(channelName, cast, testDb)
	for _, v := range dataSet {
		air, err := NewAir(v.begin, v.end)
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
	curr, err := ch.CurrentAirTime()
	if err != nil {
		t.Errorf("getting current time %v", err)
	}
	if curr.UUID != dataSet[1].uuid {
		t.Errorf("expected %v go %v", dataSet[1].uuid, curr.UUID)
	}
}

func TestGetChannel(t *testing.T) {
	ch, err := GetChannel(channelName, cast, testDb)
	if err != nil {
		t.Errorf("getting channel :%v", err)
	}
	if ch.Name != channelName {
		t.Errorf("expected %s got %s", channelName, ch.Name)
	}

	curr, err := ch.CurrentAirTime()
	if err == nil {
		t.Error("expected error got nil")
	}
	if curr != nil {
		t.Errorf("expected nil got %v", curr)
	}

	ch, err = GetChannel("", cast, testDb)
	if err == nil {
		t.Error("expected an error got nil")
	}
	if ch != nil {
		t.Errorf("expected nil got %v", ch)
	}

}

func TestChannel_Exists(t *testing.T) {
	ch, err := GetChannel(channelName, cast, testDb)
	if err != nil {
		t.Errorf("getting channel :%v", err)
	}
	e := ch.timeTable.table[0].event.(*Air)
	if !ch.Exists(e) {
		t.Errorf(" expected the airtime  to exist")
	}
	err = ch.AddAirTime(e)
	if err == nil {
		t.Error("expected error got nil")
	}
}

func TestCleanUp(t *testing.T) {
	testDb.DeleteDatabase()
}

func iter(n int) []struct{} {
	return make([]struct{}, n)
}
