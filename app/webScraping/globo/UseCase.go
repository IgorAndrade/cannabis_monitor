package globo

import (
	"fmt"

	"github.com/IgorAndrade/cannabis_monitor/app/webScraping"
	"github.com/IgorAndrade/cannabis_monitor/internal/repository"
)

const baseURL = "https://www.globo.com/busca/"
const suffixURL = "&page=%d&order=recent&from=now-1d"

type Explorer struct {
	elastic   repository.Elasticsearch
	baseURL   string
	suffixURL string
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
	}
	for _, f := range fnc {
		f(e)
	}
	return e
}

func (d Explorer) Search(words []string) {
	urls := make([]string, len(words))
	for i, w := range words {
		url := fmt.Sprintf("%s?q=%s%s", d.baseURL, w, d.suffixURL)
		urls[i] = url
	}
	ch := Scraping(urls)
	for result := range ch {
		d.elastic.Post(result)
	}
}
