package webserver

import "honeypot/timelines"

type env struct {
	tl timelines.Timelines
}

func newEnv() *env {
	return &env{
		tl: timelines.InitTimelines(),
	}
}

func (e *env) destroy() {
	e.tl.Close()
}
