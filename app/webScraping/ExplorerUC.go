package webScraping

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"time"

	"golang.org/x/sync/errgroup"
)

type Explorer interface {
	Search(words []string)
}

type QueryResult struct {
	Publisher string `json:"publisher,omitempty"`
	Page      string `json:"page,omitempty"`
	Title     string `json:"title,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Link      string `json:"link,omitempty"`
	Date      string `json:"Date,omitempty"`
}

func (q QueryResult) GetID() string {
	h := sha1.New()
	h.Write([]byte(q.Title))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

func (q QueryResult) String() string {
	if b, err := json.Marshal(q); err == nil {
		return string(b)
	}
	return ""
}

type Clocker interface {
	Now() time.Time
}

type ClockerImp struct {
	Fnc func() time.Time
}

func (c ClockerImp) Now() time.Time {
	if c.Fnc == nil {
		return time.Now().Local()
	}
	return c.Fnc()
}

type Scraping []Explorer

func (list Scraping) Search(words []string) error {
	g, _ := errgroup.WithContext(context.TODO())
	for _, s := range list {
		g.Go(func() error {
			s.Search(words)
			return nil
		})
	}
	return g.Wait()
}
