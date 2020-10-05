package estadao

import (
	"fmt"

	"github.com/IgorAndrade/cannabis_monitor/app/webScraping"
	"github.com/IgorAndrade/cannabis_monitor/internal/repository"
)

const baseURL = "https://busca.estadao.com.br/?tipo_conteudo=Todos"

//https://busca.estadao.com.br/?tipo_conteudo=Todos&quando=nas-ultimas-24-horas&q=cannabis

type Explorer struct {
	elastic repository.Elasticsearch
	BaseURL string
	When    string
	Query   string
}

type ExplorerConf func(*Explorer)

func WithWhen(when string) ExplorerConf {
	return func(e *Explorer) {
		e.When = fmt.Sprintf("&quando=%s", when)
	}
}
func NewExplorer(rep repository.Elasticsearch, fnc ...ExplorerConf) webScraping.Explorer {
	e := &Explorer{
		elastic: rep,
		BaseURL: baseURL,
		When:    "&quando=nas-ultimas-24-horas",
		Query:   "&q=cannabis",
	}
	for _, f := range fnc {
		f(e)
	}
	return e
}

func (d Explorer) Search(words []string) {

}
