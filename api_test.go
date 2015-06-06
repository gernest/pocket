package pocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPocket(t *testing.T) {
	var (
		chName = "hey"
		bName  = "yo"
	)
	defer db.DeleteDatabase()
	ad := NewAd("hello world")
	air, err := NewAir(now, now.Add(time.Second))
	if err != nil {
		t.Errorf("creating new air %v", err)
	}
	air.AddAd(ad)
	ch := NewChannel(chName, bName, db)
	ch.AddAirTime(air)
	curr, err := ch.CurrentAirTime()
	if err != nil {
		t.Errorf("getting current time %v", err)
	}
	if curr.UUID != air.UUID {
		t.Errorf("expected %v got %v", air, curr)
	}

	// case no broadcast
	home := fmt.Sprintf("/%s/%s", bName, chName)
	req, err := http.NewRequest("GET", home, nil)
	if err != nil {
		t.Errorf("creating request %v", err)
	}
	w := httptest.NewRecorder()
	Pocket(w, req)
	if !strings.Contains(w.Body.String(), errBroadCastNotFound.Error()) {
		t.Errorf("expected %v got %s", errBroadCastNotFound, w.Body.String())
	}

	// case there broadcast but no registered channel
	b := NewBroadcast(bName, db)
	err = b.Save()
	if err != nil {
		t.Errorf("creating broadcast %v", err)
	}
	w2 := httptest.NewRecorder()
	Pocket(w2, req)
	if !strings.Contains(w2.Body.String(), errChannelNotFound.Error()) {
		t.Errorf("expected %v got %s", errChannelNotFound, w2.Body.String())
	}
	b.AddChannel(ch)
	b.Save()
	w3 := httptest.NewRecorder()
	Pocket(w3, req)
	if !strings.Contains(w3.Body.String(), air.UUID) {
		t.Errorf("expected %v got %s", air, w3.Body.String())
	}

}
