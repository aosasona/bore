package aggregate

import (
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
	"go.trulyao.dev/bore/v2/pkg/errs"
)

//go:generate go tool github.com/abice/go-enum --marshal

var (
	ErrInvalidAggregateString = errs.New("invalid aggregate string")
	ErrInvalidAggregateID     = errs.New("invalid aggregate ID")
)

// An Aggregate represents a target entity for events, identified by a type and a ULID.
// This is ideally another table in the database, and are represented as strings in the format "type:id".
type Aggregate struct {
	id ulid.ULID     // valid ULID
	t  AggregateType // valid AggregateType
}

// ENUM(item,collection)
type AggregateType string

// New creates a new aggregate with the given type and a newly generated ULID.
func New(aggregateType AggregateType) (Aggregate, error) {
	if !aggregateType.IsValid() {
		return Aggregate{}, ErrInvalidAggregateString
	}

	return Aggregate{id: ulid.Make(), t: aggregateType}, nil
}

// WithID creates a new aggregate with the given type and ULID.
func WithID(aggregateType AggregateType, aggregateId string) (Aggregate, error) {
	if !aggregateType.IsValid() {
		return Aggregate{}, ErrInvalidAggregateString
	}

	id, err := ulid.Parse(strings.TrimSpace(aggregateId))
	if err != nil {
		return Aggregate{}, ErrInvalidAggregateID
	}

	if id.Compare(ulid.ULID{}) == 0 {
		return Aggregate{}, ErrInvalidAggregateID
	}

	return Aggregate{id: id, t: aggregateType}, nil
}

// FromRaw creates a new aggregate from raw string inputs for type and ID.
func FromRaw(aggregateType string, aggregateId string) (Aggregate, error) {
	return WithID(AggregateType(strings.TrimSpace(aggregateType)), aggregateId)
}

// Parses an aggregate from its string representation in the format "type:id".
func Parse(s string) (Aggregate, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Aggregate{}, ErrInvalidAggregateString
	}

	aggregateType := AggregateType(strings.TrimSpace(parts[0]))
	if !aggregateType.IsValid() {
		return Aggregate{}, ErrInvalidAggregateString
	}

	id, err := ulid.Parse(strings.TrimSpace(parts[1]))
	if err != nil {
		return Aggregate{}, ErrInvalidAggregateID
	}

	return Aggregate{id: id, t: aggregateType}, nil
}

func (a *Aggregate) IsValid() bool {
	return a.t.IsValid() && a.id.Compare(ulid.ULID{}) != 0
}

func (a *Aggregate) RawID() ulid.ULID       { return a.id }
func (a *Aggregate) RawType() AggregateType { return a.t }

func (a *Aggregate) ID() string   { return a.id.String() }
func (a *Aggregate) Type() string { return a.t.String() }

// String returns the string representation of the aggregate in the format "type:id".
func (a *Aggregate) String() string {
	return fmt.Sprintf("%s:%s", a.t, a.id)
}
