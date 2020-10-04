package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/IgorAndrade/analytics-twitter/server/app/config"
	"github.com/IgorAndrade/go-boilerplate/internal/model"
	"github.com/IgorAndrade/go-boilerplate/internal/repository"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mitchellh/mapstructure"
	"github.com/yalp/jsonpath"
)

type Elasticsearch struct {
	client       *elasticsearch.Client
	index        string
	documenttype string
}

func newServer(cfg config.Elasticsearch, index string, documenttype string) (repository.Elasticsearch, error) {
	elsCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.Address,
		},
		Username: cfg.Username,
		Password: cfg.Password,
	}
	client, err := elasticsearch.NewClient(elsCfg)
	if err != nil {
		return nil, err
	}
	// Test connect
	res, err := client.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	log.Println(res)

	return &Elasticsearch{
		client:       client,
		documenttype: documenttype,
		index:        index,
	}, nil
}

func (s Elasticsearch) Post(p model.Post) error {
	req := esapi.IndexRequest{
		Index:        s.index,
		DocumentType: s.documenttype,
		DocumentID:   p.GetID(),
		Body:         strings.NewReader(p.String()),
		Refresh:      "true",
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	res, err := req.Do(ctx, s.client)
	if res != nil {
		res.Body.Close()
	}
	fmt.Println(p)
	return err
}

func (s Elasticsearch) Find(ctx context.Context, query map[string]string) ([]model.Post, error) {

	buf := new(bytes.Buffer)
	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query,
		},
	}
	json.NewEncoder(buf).Encode(queryBody)
	es, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	return adpter(es.Body)
}

func adpter(rc io.ReadCloser) ([]model.Post, error) {
	var posts []model.Post
	var data interface{}
	if err := json.NewDecoder(rc).Decode(&data); err != nil {
		return posts, err
	}
	raw, err := jsonpath.Read(data, "$.._source")
	if err != nil {
		return posts, err
	}

	mapstructure.Decode(raw, &posts)
	return posts, nil
}
