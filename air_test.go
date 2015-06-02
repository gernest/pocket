package pocket

import (
	"testing"
	"time"
)

func TestNewAir(t *testing.T) {
	begin := time.Now()
	end := begin.Add(time.Minute)

	a1, err := NewAIR(begin, end)
	if err != nil {
		t.Errorf("creating new air: %v", err)
	}
	if !a1.Begin.Equal(begin) {
		t.Errorf("expected %s got %s", begin, a1.Begin)
	}
	a2, err := NewAIR(begin.String(), end.String())
	if err == nil {
		t.Errorf("expected an error %v", err)
	}
	if a2 != nil {
		t.Errorf("expected nil got %v", a2)
	}
	a3, err := NewAIR(begin.Format(time.RFC822), end.Format(time.RFC822))
	if err != nil {
		t.Errorf("creating new air: %v", err)
	}
	if a3.Begin.Format(time.RFC822) != begin.Format(time.RFC822) {
		t.Errorf("expected %s got %s", begin.Format(time.RFC822), a3.Begin.Format(time.RFC822))
	}
	a4, err := NewAIR(begin, end.String())
	if err == nil {
		t.Errorf("expected an error %v", err)
	}
	if a4 != nil {
		t.Errorf("expected nil got %v", a2)
	}
}
