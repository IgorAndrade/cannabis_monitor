package estadao

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

//const baseURL = "https://busca.estadao.com.br/?tipo_conteudo=Todos"
const baseURL = "https://busca.estadao.com.br/modulos/busca-resultado?modulo=busca-resultado&config%5Bbusca%5D%5Bpage%5D=0&config%5Bbusca%5D%5Bparams%5D=tipo_conteudo%3DTodos%26quando%3Dno-ultimo-mes%26q%3DWORD&ajax=1"
const PAGE = "&config%5Bbusca%5D%5Bpage%5D="
const PUBLISHER = "Estadao"

//https://busca.estadao.com.br/modulos/busca-resultado?modulo=busca-resultado&config[busca][page]=3&config[busca][params]=tipo_conteudo=Todos&quando=no-ultimo-mes&q=maconha&ajax=1
//https://busca.estadao.com.br/?tipo_conteudo=Todos&quando=nas-ultimas-24-horas&q=cannabis

type Explorer struct {
	elastic repository.Elasticsearch
	BaseURL string
	When    string
	Clocker webScraping.Clocker
}

type ExplorerConf func(*Explorer)

func WithWhen(when string) ExplorerConf {
	return func(e *Explorer) {

		e.BaseURL = strings.Replace(e.BaseURL, "no-ultimo-mes", when, 1)
	}
}
func NewExplorer(rep repository.Elasticsearch, fnc ...ExplorerConf) webScraping.Explorer {
	e := &Explorer{
		elastic: rep,
		BaseURL: baseURL,
		Clocker: webScraping.ClockerImp{},
	}
	for _, f := range fnc {
		f(e)
	}
	return e
}

func (e Explorer) Search(words []string) {
	if yml.String("Estadao.Enable", "false") != "true" {
		return
	}
	urls := make([]string, len(words))
	for i, w := range words {
		url := strings.Replace(e.BaseURL, "WORD", w, 1)
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
		for p := 1; p < 30; p++ {
			urlPage := strings.Replace(url, "page%5D=0", "page%5D="+strconv.Itoa(p), 1)
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

			if document.Find(".item-lista-busca").Size() < 1 {
				return nil
			}
			document.Find(".item-lista-busca").Each(func(i int, s *goquery.Selection) {
				p := s.Find(".cor-e").Text()

				contA := s.Find(".link-title")
				a, _ := contA.Attr("href")
				t := contA.Text()

				d := contA.Find("p").Text()
				strDate := s.Find(".data-posts").Text()
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
	var re = regexp.MustCompile(`(\d\d) de\s+(\S+)\s+de\s+(\d{4}).+\s(\d{1,2})h(\d{1,2})`)
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
	"janeiro":   "01",
	"fevereiro": "02",
	"março":     "03",
	"abril":     "04",
	"maio":      "05",
	"junho":     "06",
	"julho":     "07",
	"agosto":    "08",
	"setembro":  "09",
	"outubro":   "10",
	"novembro":  "11",
	"dezembro":  "12",
}
