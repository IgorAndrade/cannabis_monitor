package estadao

import (
	"fmt"
	"testing"

	"github.com/IgorAndrade/cannabis_monitor/app/webScraping"
	"github.com/IgorAndrade/cannabis_monitor/internal/repository"
)

func TestExplorer_Scraping(t *testing.T) {
	type fields struct {
		elastic repository.Elasticsearch
		BaseURL string
		When    string
		Clocker webScraping.Clocker
	}
	type args struct {
		urls []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan webScraping.QueryResult
	}{
		{
			name: "teste url",
			args: args{
				urls: []string{"https://busca.estadao.com.br/?tipo_conteudo=Todos&quando=no-ultimo-mes&q=cannabis"},
			},
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Explorer{}
			ch := e.Scraping(tt.args.urls)
			for msg := range ch {
				fmt.Println(msg)
			}
		})
	}
}
