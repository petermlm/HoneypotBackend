package webserver

import "honeypot/timelines"

type env struct {
	tl timelines.TimelinesQuery
}

func newEnv() *env {
	return &env{
		tl: timelines.NewTimelinesQuery(),
	}
}

func (e *env) destroy() {
	e.tl.Close()
}
