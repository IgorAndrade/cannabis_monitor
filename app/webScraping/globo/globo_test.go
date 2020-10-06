package globo

import (
	"fmt"
	"testing"
	"time"

	"github.com/IgorAndrade/cannabis_monitor/app/webScraping"
)

const DefaultTime = "2020-10-05T21:13:12-03:00"

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
			e := &Explorer{}
			ch := e.Scraping(tt.args.urls)
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
			name: "default",
			args: args{
				str: "há 4 default",
			},
			want: "2020-10-05T21:13:12-03:00",
		},
		{
			name: "dias",
			args: args{
				str: "há 4 dias",
			},
			want: "2020-10-01T21:13:12-03:00",
		},
		{
			name: "horas",
			args: args{
				str: "há 10 horas",
			},
			want: "2020-10-05T11:13:12-03:00",
		},
		{
			name: "minutos",
			args: args{
				str: "há 12 minutos",
			},
			want: "2020-10-05T21:01:12-03:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Explorer{}
			e.Clocker = webScraping.ClockerImp{Fnc: func() time.Time {
				defaultTime, err := time.Parse(time.RFC3339, DefaultTime)
				if err != nil {
					t.Error(err)
				}
				return defaultTime
			}}

			if got := e.parseDate(tt.args.str); got != tt.want {
				t.Errorf("parseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
