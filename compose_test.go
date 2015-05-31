package pocket

import (
	"strings"
	"testing"
	"time"
)

func TestComposer_At(t *testing.T) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	c := NewComposer().At(tomorrow)
	if !strings.Contains(c.String(), every) {
		t.Errorf("expected %s to contain %s", c, every)
	}
}
