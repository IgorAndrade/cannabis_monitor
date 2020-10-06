package globo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IgorAndrade/cannabis_monitor/app/webScraping"
	"github.com/IgorAndrade/cannabis_monitor/internal/repository"
	"github.com/PuerkitoBio/goquery"
	yml "github.com/gookit/config/v2"
	"golang.org/x/sync/errgroup"
)

const baseURL = "https://www.globo.com/busca/"
const suffixURL = "&page=%d&order=recent&from=now-1d"
const PUBLISHER = "Globo"

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

func WithSuffixURL(suffixURL string) ExplorerConf {
	return func(e *Explorer) {
		e.suffixURL = suffixURL
	}
}

func NewExplorer(rep repository.Elasticsearch, fnc ...ExplorerConf) webScraping.Explorer {
	e := &Explorer{
		elastic:   rep,
		baseURL:   baseURL,
		suffixURL: suffixURL,
		Clocker:   webScraping.ClockerImp{},
	}
	for _, f := range fnc {
		f(e)
	}
	return e
}

func (e Explorer) Search(words []string) {
	if yml.String("Globo.Enable", "false") != "true" {
		return
	}
	urls := make([]string, len(words))
	for i, w := range words {
		url := fmt.Sprintf("%s?q=%s%s", e.baseURL, w, e.suffixURL)
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
	}
}

func (e Explorer) parseDate(str string) string {
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
		return e.Clocker.Now().Add(time.Duration(-1*num) * time.Hour).Format(time.RFC3339)
	}
	if strings.Contains(str, "dias") {
		return e.Clocker.Now().Add(time.Duration(-24*num) * time.Hour).Format(time.RFC3339)
	}
	if strings.Contains(str, "dia") {
		return e.Clocker.Now().Add(time.Duration(-24*num) * time.Hour).Format(time.RFC3339)
	}
	if strings.Contains(str, "minutos") {
		return e.Clocker.Now().Add(time.Duration(-1*num) * time.Minute).Format(time.RFC3339)
	}
	return e.Clocker.Now().Format(time.RFC3339)
}
