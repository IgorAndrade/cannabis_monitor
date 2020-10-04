package globo

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

type QueryResult struct {
	Page   string `json:"page,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Link   string `json:"link,omitempty"`
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

func Scraping(urls []string) <-chan QueryResult {
	ch := make(chan QueryResult, 10)
	g, _ := errgroup.WithContext(context.TODO())
	for _, url := range urls {
		g.Go(func() error {
			return scraping(url, ch)
		})
	}
	go func(g *errgroup.Group, ch chan QueryResult) {
		g.Wait()
		close(ch)
	}(g, ch)
	return ch
}

func scraping(url string, ch chan QueryResult) error {
	for p := 1; ; p++ {
		urlPage := fmt.Sprintf(url, p)
		fmt.Printf("url: %s \n\n", urlPage)
		response, err := http.Get(urlPage)
		if err != nil {
			log.Println(err)
			return err
		}
		defer response.Body.Close()

		// Create a goquery document from the HTTP response
		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Printf("Error loading HTTP response body. %s \n", err.Error())
			return err
		}

		if document.Find(".widget--info__text-container").Size() < 1 {
			return nil
		}
		document.Find(".widget--info__text-container").Each(func(i int, s *goquery.Selection) {
			p := s.ContentsFiltered(".widget--info__header").Text()
			contA := s.ContentsFiltered("a")
			a, _ := contA.Attr("href")
			t := contA.ContentsFiltered(".widget--info__title").Text()
			d := contA.ContentsFiltered(".widget--info__description").Text()
			q := QueryResult{
				Page:   p,
				Detail: strings.TrimSpace(d),
				Title:  strings.TrimSpace(t),
				Link:   a,
			}
			fmt.Println(q)
			ch <- q
		})
	}
}
