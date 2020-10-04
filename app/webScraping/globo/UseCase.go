package globo

import (
	"fmt"

	"github.com/IgorAndrade/go-boilerplate/internal/repository"
)

const baseURL = "https://www.globo.com/busca/"
const suffixURL = "&page=%d&order=recent&from=now-1d"

type Discover struct {
	elastic repository.Elasticsearch
}

func (d Discover) search(words []string) {
	urls := make([]string, len(words))
	for i, w := range words {
		url := fmt.Sprintf("%s?q=%s,%s", baseURL, w, suffixURL)
		urls[i] = url
	}
	ch := Scraping(urls)
	for result := range ch {
		d.elastic.Post(result)
	}
}
