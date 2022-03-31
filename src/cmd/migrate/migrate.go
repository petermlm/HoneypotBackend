package main

import (
	"honeypot/timelines"
)

func main() {
	tl := timelines.InitTimelines()

	if !tl.MigrationsTableExsits() {
		tl.MigrateCmd([]string{"init"})
	}

	tl.MigrateCmd([]string{})
}
