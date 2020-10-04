package repository

import (
	"github.com/IgorAndrade/go-boilerplate/internal/model"
)

type Elasticsearch interface {
	Post(m model.Post) error
}
