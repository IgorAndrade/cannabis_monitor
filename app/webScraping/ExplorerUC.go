package webScraping

import "time"

type Explorer interface {
	Search(words []string)
}

type Clocker interface {
	Now() time.Time
}

type ClockerImp struct {
	Fnc func() time.Time
}

func (c ClockerImp) Now() time.Time {
	if c.Fnc == nil {
		return time.Now()
	}
	return c.Fnc()
}
