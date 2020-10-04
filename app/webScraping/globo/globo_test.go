package globo

import (
	"fmt"
	"testing"
)

func TestScraping(t *testing.T) {
	type args struct {
		urls []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "readind globo",
			args: args{
				urls: []string{
					"https://www.globo.com/busca/?q=maconha&page=%d&order=recent&from=now-1d",
					"https://www.globo.com/busca/?q=cannabis&page=%d&order=recent&from=now-1d",
					"https://www.globo.com/busca/?q=maconheiro&page=%d&order=recent&from=now-1d",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := Scraping(tt.args.urls)
			for msg := range ch {
				fmt.Println(msg)
			}
		})
	}
}
