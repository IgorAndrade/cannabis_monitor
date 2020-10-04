package globo

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

type QueryResult struct {
	Page   string `json:"page,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Link   string `json:"link,omitempty"`
	Date   string `json:"Date,omitempty"`
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
		g.Go(scraping(url, ch))
	}
	go func(g *errgroup.Group, ch chan QueryResult) {
		g.Wait()
		close(ch)
	}(g, ch)
	return ch
}

func scraping(url string, ch chan QueryResult) func() error {
	return func() error {
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
				strDate := contA.ContentsFiltered(".widget--info__meta").Text()
				q := QueryResult{
					Page:   p,
					Detail: strings.TrimSpace(d),
					Title:  strings.TrimSpace(t),
					Link:   a,
					Date:   parseDate(strDate),
				}
				//fmt.Println(q)
				ch <- q
			})
		}
	}
}

func parseDate(str string) string {
	var re = regexp.MustCompile(`\d+`)
	var num int = 1
	for _, match := range re.FindAllString(str, -1) {
		n, err := strconv.Atoi(match)
		if err != nil {
			return time.Now().Format(time.RFC3339)
		}
		num = n
		break
	}
	if strings.Contains(str, "horas") {
		return time.Now().Add(time.Duration(-1*num) * time.Hour).Format(time.RFC3339)
	}
	if strings.Contains(str, "dias") {
		return time.Now().Add(time.Duration(-24*num) * time.Hour).Format(time.RFC3339)
	}
	if strings.Contains(str, "dia") {
		return time.Now().Add(time.Duration(-24*num) * time.Hour).Format(time.RFC3339)
	}
	if strings.Contains(str, "minutos") {
		return time.Now().Add(time.Duration(-1*num) * time.Minute).Format(time.RFC3339)
	}
	return time.Now().Format(time.RFC3339)
}
