package aggregate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

//go:generate go tool github.com/abice/go-enum --marshal

var (
	ErrInvalidAggregateString = errors.New("invalid aggregate string")
	ErrInvalidAggregateID     = errors.New("invalid aggregate ID")
)

// An aggregate represents a target entity for events, identified by a type and a ULID.
// This is ideally another table in the database, and are represented as strings in the format "type:id".
type aggregate struct {
	id ulid.ULID     // valid ULID
	t  AggregateType // valid AggregateType
}

// ENUM(item,collection)
type AggregateType string

// Creates a new aggregate with the given type and ULID.
func New(aggregateType AggregateType, id ulid.ULID) (*aggregate, error) {
	if !aggregateType.IsValid() {
		return nil, ErrInvalidAggregateString
	}

	if id.Compare(ulid.ULID{}) == 0 {
		return nil, ErrInvalidAggregateID
	}

	return &aggregate{id: id, t: aggregateType}, nil
}

// Parses an aggregate from its string representation in the format "type:id".
func Parse(s string) (aggregate, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return aggregate{}, ErrInvalidAggregateString
	}

	aggregateType := AggregateType(strings.TrimSpace(parts[0]))
	if !aggregateType.IsValid() {
		return aggregate{}, ErrInvalidAggregateString
	}

	id, err := ulid.Parse(strings.TrimSpace(parts[1]))
	if err != nil {
		return aggregate{}, ErrInvalidAggregateID
	}

	return aggregate{id: id, t: aggregateType}, nil
}

func (a *aggregate) ID() string {
	return a.id.String()
}

func (a *aggregate) Type() AggregateType {
	return a.t
}

// String returns the string representation of the aggregate in the format "type:id".
func (a *aggregate) String() string {
	return fmt.Sprintf("%s:%s", a.t, a.id)
}
