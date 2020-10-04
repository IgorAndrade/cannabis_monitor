package main

import (
	"log"

	"github.com/IgorAndrade/cannabis_monitor/app/api/elasticsearch"
	"github.com/IgorAndrade/cannabis_monitor/app/config"
	"github.com/IgorAndrade/cannabis_monitor/app/webScraping/globo"
)

func main() {
	conf := config.GetConfi()
	elsConf := conf.Elasticsearch
	globoConf := conf.Globo
	rep, err := elasticsearch.NewServer(elsConf, globoConf.Index, globoConf.DocumentType)
	if err != nil {
		log.Fatalln(err)
	}
	serv := globo.NewExplorer(rep)
	serv.Search([]string{"maconha", "cannabis", "legalização"})
}
