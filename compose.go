package pocket

import (
	"bytes"
	"time"
)

var (
	every = "@every"
	space = " "
)

// Composer composes cron jobs strings.
type Composer struct {
	buf *bytes.Buffer
}

func NewComposer() *Composer {
	return &Composer{
		buf: &bytes.Buffer{},
	}
}

func (c *Composer) At(t time.Time) *Composer {
	d := t.Sub(time.Now())
	c.buf.WriteString(every + space + d.String())
	return c
}

func (c *Composer) String() string {
	return c.buf.String()
}
