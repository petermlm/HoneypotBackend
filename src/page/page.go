package page

import (
	"honeypot/settings"
)

type Page struct {
	Num  int
	Size int
}

func NewPage() Page {
	return Page{
		Num:  settings.PageDefaultNum,
		Size: settings.PageDefaultSize,
	}
}

func (p Page) Limit() int {
	return p.Size
}

func (p Page) Offset() int {
	return p.Num * p.Size
}
