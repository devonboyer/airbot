package airtable

import "context"

type Interface interface {
	Base(baseID string) BaseInterface
}

type BaseInterface interface {
	Table(name string) TableInterface
}

type TableInterface interface {
	List(ctx context.Context, opts *ListOptions, v interface{}) error
}
