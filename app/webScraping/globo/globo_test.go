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

func Test_parseDate(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "dias",
			args: args{
				str: "há 4 dias",
			},
			want: "",
		},
		{
			name: "horas",
			args: args{
				str: "há 10 horas",
			},
			want: "",
		},
		{
			name: "minutos",
			args: args{
				str: "há 18 minutos",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDate(tt.args.str); got != tt.want {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
