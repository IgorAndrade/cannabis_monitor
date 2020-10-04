package model

import "fmt"

type Post interface {
	fmt.Stringer
	GetID() string
}
