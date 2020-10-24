package folha

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/IgorAndrade/cannabis_monitor/app/webScraping"
	"github.com/IgorAndrade/cannabis_monitor/internal/repository"
	"github.com/PuerkitoBio/goquery"
	yml "github.com/gookit/config/v2"
	"golang.org/x/sync/errgroup"
)

const baseURL = "https://search.folha.uol.com.br/search?q=WORD&site=todos&periodo=WHEN"
const PUBLISHER = "Folha"

type Explorer struct {
	elastic   repository.Elasticsearch
	baseURL   string
	suffixURL string
	Clocker   webScraping.Clocker
}

type ExplorerConf func(*Explorer)

func WithBaseURL(baseURL string) ExplorerConf {
	return func(e *Explorer) {
		e.baseURL = baseURL
	}
}
func WithWhen(when string) ExplorerConf {
	return func(e *Explorer) {
		e.baseURL = strings.Replace(e.baseURL, "WHEN", when, 1)
	}
}
func NewExplorer(rep repository.Elasticsearch, fnc ...ExplorerConf) webScraping.Explorer {
	e := &Explorer{
		elastic: rep,
		baseURL: strings.Replace(baseURL, "WHEN", "mes", 1),
		Clocker: webScraping.ClockerImp{},
	}
	for _, f := range fnc {
		f(e)
	}
	return e
}

func (e Explorer) Search(words []string) {
	if yml.String("Folha.Enable", "false") != "true" {
		return
	}
	urls := make([]string, len(words))
	for i, w := range words {
		url := strings.Replace(e.baseURL, "WORD", w, 1)
		urls[i] = url
	}
	ch := e.Scraping(urls)
	for result := range ch {
		e.elastic.Post(result)
	}
}

func (e Explorer) Scraping(urls []string) <-chan webScraping.QueryResult {
	ch := make(chan webScraping.QueryResult, 10)
	g, _ := errgroup.WithContext(context.TODO())
	for _, url := range urls {
		g.Go(e.scraping(url, ch))
	}
	go func(g *errgroup.Group, ch chan webScraping.QueryResult) {
		g.Wait()
		close(ch)
	}(g, ch)
	return ch
}

func (e Explorer) scraping(url string, ch chan webScraping.QueryResult) func() error {
	return func() error {
		for p := 1; p < 30; p += 10 {
			urlPage := fmt.Sprintf("%s&sr=%d", url, p)
			fmt.Printf("url: %s \n\n page %d \n\n", urlPage, p)
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

			if document.Find(".c-headline--newslist").Size() < 1 {
				return nil
			}
			document.Find(".c-headline--newslist").Each(func(i int, s *goquery.Selection) {
				p := s.Find(".c-search__result_h3").Text()

				contA := s.Find(".c-headline__content a")
				a, _ := contA.Attr("href")
				t := contA.Find(".c-headline__title").Text()

				d := contA.Find("p").Text()
				strDate := s.Find("time").Text()
				q := webScraping.QueryResult{
					Publisher: PUBLISHER,
					Page:      p,
					Detail:    strings.TrimSpace(d),
					Title:     strings.TrimSpace(t),
					Link:      a,
					Date:      e.parseDate(strDate),
				}
				ch <- q
			})
		}
		return nil
	}
}

func (e Explorer) parseDate(str string) string {
	var re = regexp.MustCompile(`(\d{1,2})\.(\S{3})\.(\d{4})\D+(\d{1,2})h(\d{1,2})`)
	match := re.FindAllStringSubmatch(str, -1)
	if len(match) == 0 {
		return e.Clocker.Now().Format(time.RFC3339)
	}
	day := match[0][1]
	mouth := match[0][2]
	year := match[0][3]
	h := match[0][4]
	m := match[0][5]
	date := fmt.Sprintf("%s/%s/%s %s:%s -0300", day, mapDate[mouth], year, h, m)
	t, _ := time.Parse("02/01/2006 03:04 -0700", date)
	return t.Format(time.RFC3339)
}

var mapDate = map[string]string{
	"jan": "01",
	"fev": "02",
	"mar": "03",
	"abr": "04",
	"mai": "05",
	"jun": "06",
	"jul": "07",
	"ago": "08",
	"set": "09",
	"out": "10",
	"nov": "11",
	"dez": "12",
}
