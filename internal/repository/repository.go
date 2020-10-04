package repository

import (
	"github.com/IgorAndrade/cannabis_monitor/internal/model"
)

type Elasticsearch interface {
	Post(m model.Post) error
}
