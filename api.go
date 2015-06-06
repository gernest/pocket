package pocket

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gernest/nutz"
)

var (
	db                   = nutz.NewStorage("pocket.db", 0600, nil)
	errChannelNotFound   = errors.New("channel not found")
	errBroadCastNotFound = errors.New("broadcast not found")
	noAir                = "no air"
)

func Pocket(w http.ResponseWriter, r *http.Request) {
	var (
		slash  = "/"
		result = func(key string, v interface{}) map[string]interface{} {
			m := make(map[string]interface{})
			m[key] = v
			return m
		}
		jErr = func(err error) map[string]interface{} {
			return result("error", err.Error())
		}
	)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := json.NewEncoder(w)

	p := strings.Trim(r.URL.Path, slash)
	s := strings.Split(p, slash)
	if len(s) == 2 {
		cast := s[0]
		channel := s[1]
		bd, err := GetBroadCast(cast, db)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			e.Encode(jErr(errBroadCastNotFound))
			return
		}
		if !bd.HasChannel(channel) {
			w.WriteHeader(http.StatusOK)
			e.Encode(jErr(errChannelNotFound))
			return
		}
		ch, err := GetChannel(channel, cast, db)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			e.Encode(jErr(errChannelNotFound))
			return
		}
		if r.Method == "GET" {
			on, err := ch.CurrentAirTime()
			if err != nil {
				if err.Error() == errScheduleNotFound.Error() {
					w.WriteHeader(http.StatusOK)
					e.Encode(result("noAir", noAir))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				e.Encode(jErr(err))
				return
			}
			w.WriteHeader(http.StatusOK)
			e.Encode(result("air", on))
			return
		}

	}
}
