package main

import (
	"context"
	"flag"
	"log"

	"github.com/IgorAndrade/cannabis_monitor/app/api/elasticsearch"
	"github.com/IgorAndrade/cannabis_monitor/app/config"
	"github.com/IgorAndrade/cannabis_monitor/app/webScraping/estadao"
	"github.com/IgorAndrade/cannabis_monitor/app/webScraping/globo"
	yml "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"golang.org/x/sync/errgroup"
)

type Ya struct {
	Name    string
	Debug   bool
	BaseKey string
	Other   other
}

type other struct {
	Name  string
	Debug bool
}

func main() {
	file := flag.String("config", "./config.yaml", " -config=config.yaml")
	flag.Parse()
	yml.AddDriver(yaml.Driver)
	yml.LoadFiles(*file)
	var elsConf config.Elasticsearch
	yml.BindStruct("Elasticsearch", &elsConf)

	rep, err := elasticsearch.NewServer(elsConf)
	if err != nil {
		log.Fatalln(err)
	}

	words := yml.Strings("words")
	globo := globo.NewExplorer(rep)
	estadao := estadao.NewExplorer(rep)
	g, _ := errgroup.WithContext(context.TODO())
	g.Go(func() error {
		globo.Search(words)
		return nil
	})
	g.Go(func() error {
		estadao.Search(words)
		return nil
	})

	g.Wait()
}
