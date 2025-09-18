package events

import "go.trulyao.dev/bore/v2/pkg/events/aggregate"

//go:generate go tool github.com/abice/go-enum --marshal

// ENUM(create_item,delete_item,create_collection,delete_collection)
type Action string

type Event struct {
	ID        string
	Aggregate aggregate.Aggregate
}
