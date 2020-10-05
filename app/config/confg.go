package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/sarulabs/di"
)

type Elasticsearch struct {
	Address  string
	Username string
	Password string
}

type Globo struct {
	Index        string
	DocumentType string
}
type Config struct {
	Rest          Rest
	Mongo         Mongo
	Elasticsearch Elasticsearch
}

type Rest struct {
	Port string
}

type Mongo struct {
	Address  string
	User     string
	Password string
}

var CONFIG = "config"
var c *Config
var once = sync.Once{}

func GetConfi() Config {
	once.Do(func() {
		c = &Config{
			Rest: Rest{
				Port: fmt.Sprintf(":%s", os.Getenv("API_PORT")),
			},
			Mongo: Mongo{
				Address:  os.Getenv("MONGO_URL"),
				User:     os.Getenv("MONGO_USER"),
				Password: os.Getenv("MONGO_PASSWORD"),
			},
			Elasticsearch: Elasticsearch{
				Address:  os.Getenv("ELASTICSEARCH_ADDRESS"),
				Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
				Username: os.Getenv("ELASTICSEARCH_USERNAME"),
			},
		}
	})
	return *c
}
func Define(b *di.Builder) {
	b.Add(di.Def{
		Name:  CONFIG,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return GetConfi(), nil
		},
	})
}
