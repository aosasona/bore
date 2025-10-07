package events

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
	"go.trulyao.dev/bore/v2/pkg/events/payload"
)

var (
	ErrUnknownEventType = errs.New("unknown event type")
	ErrNoEventsToAppend = errs.New("no events to append")
	ErrNoEventApplied   = errs.New("no event applied")
)

// Manager handles event sourcing operations.
type Manager struct {
	db   *bun.DB
	repo repository.Repository
}

type AppendOptions struct {
	ExpectedVersion int64 // If equal or greater than zero, the current aggregate version must match this value.
}

func DefaultAppendOptions() AppendOptions {
	return AppendOptions{ExpectedVersion: -1}
}

func NewManager(db *bun.DB, repo repository.Repository) *Manager {
	return &Manager{db: db, repo: repo}
}

// ApplyN appends the given batch of events to the SAME aggregate, figures out the appropriate versioning, and applies the projections in a shared transaction.
func (m *Manager) ApplyN(
	ctx context.Context,
	agg aggregate.Aggregate,
	events []Event,
	options AppendOptions,
) (persisted []Event, newVersion int64, err error) {
	if len(events) == 0 {
		return nil, 0, ErrNoEventsToAppend
	} else if !agg.IsValid() {
		return nil, 0, ErrInvalidAggregate
	}

	err = m.db.RunInTx(
		ctx,
		&sql.TxOptions{Isolation: 0, ReadOnly: false},
		func(ctx context.Context, tx bun.Tx) error {
			var currentVersion int64
			err := tx.NewSelect().
				Table("events").
				ColumnExpr("IFNULL(MAX(aggregate_version), 0)").
				Where("aggregate_type = ? AND aggregate_id = ?", agg.Type(), agg.ID()).
				Scan(ctx, &currentVersion)
			if err != nil {
				return err
			}

			if options.ExpectedVersion >= 0 && currentVersion != options.ExpectedVersion {
				return errs.New("aggregate version mismatch")
			}

			timestamp := time.Now().UTC()

			rows := make([]*Event, 0, len(events))
			for i := range events {
				event := &events[i]
				event.AggregateVersion = currentVersion + int64(i) + 1

				if err := event.SetAggregate(agg); err != nil {
					return err
				}

				if event.OccurredAt.IsZero() {
					event.OccurredAt = timestamp
				}

				rows = append(rows, event)
			}

			if _, err := tx.NewInsert().Model(&rows).Ignore().Exec(ctx); err != nil {
				return err
			}

			// We need to re-fetch the events to get the auto-generated fields (like Sequence).
			var savedEvents []Event
			err = tx.NewSelect().
				Model(&savedEvents).
				Where("aggregate_type = ? AND aggregate_id = ?", agg.Type(), agg.ID()).
				Where("aggregate_version > ?", currentVersion).
				Order("aggregate_version ASC").
				Scan(ctx)
			if err != nil {
				return err
			}

			for i := range savedEvents {
				event := &savedEvents[i]
				if err := m.applyProjection(ctx, tx, event); err != nil {
					return err
				}
			}

			persisted = savedEvents
			newVersion = currentVersion + int64(len(savedEvents))

			return nil
		},
	)

	return persisted, newVersion, err
}

// Apply appends a single event to the event store and applies the projection in a shared transaction.
func (m *Manager) Apply(
	ctx context.Context,
	event *Event,
	options ...AppendOptions,
) (Event, int64, error) {
	opts := DefaultAppendOptions()
	if len(options) > 0 {
		opts = options[0]
	}

	events, newVersion, err := m.ApplyN(ctx, event.Aggregate, []Event{*event}, opts)
	if err != nil {
		return Event{}, 0, err
	}

	if len(events) == 0 {
		return Event{}, 0, ErrNoEventApplied
	}

	return events[0], newVersion, nil
}

func (m *Manager) applyProjection(ctx context.Context, tx bun.Tx, event *Event) error {
	if event == nil {
		return errs.New("event cannot be nil")
	}

	p, err := payload.Decode(event.Payload, event.Type)
	if err != nil {
		return err
	}

	options := payload.ProjectionOptions{
		Aggregate: event.Aggregate,
		Sequence:  event.Sequence,
	}

	return p.ApplyProjection(ctx, tx, m.repo, options)
}
