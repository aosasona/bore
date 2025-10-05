package bore

import (
	"context"
	"strings"

	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
	"go.trulyao.dev/bore/v2/pkg/events/payload"
	"go.trulyao.dev/bore/v2/pkg/lib"
	"go.trulyao.dev/bore/v2/pkg/mimetype"
)

type clipboardNamespace struct {
	*Bore
}

type (
	SetClipboardOptions struct {
		Passthrough  bool   // Whether to also copy to the system clipboard if available.
		CollectionID string // Optional collection ID to associate with the copied item.
		Mimetype     mimetype.MimeType
	}

	GetClipboardOptions struct {
		ItemID              string // Optional item identifier to filter pasted items.
		CollectionID        string // Optional collection ID to filter pasted items.
		FromSystemClipboard bool   // Whether to paste from the system clipboard if available.
		DeleteAfterPaste    bool   // Whether to delete the pasted item after pasting.
		SkipCollectionCheck bool   // Whether to skip checking if the collection exists.
	}

	PasteResult struct {
		Content []byte
		Item    *models.Item
	}
)

// Set copies the provided data to the Bore instance.
func (i *clipboardNamespace) Set(ctx context.Context, data []byte, opts SetClipboardOptions) error {
	forwardToSystemClipboard := i.config.ClipboardPassthrough || opts.Passthrough
	if i.clipboard.Available() && forwardToSystemClipboard {
		if err := i.clipboard.Write(ctx, data); err != nil {
			return err
		}
	}

	hash := lib.ComputeChecksum(data)
	existingItem, err := i.repository.Items().FindByHash(ctx, hash, opts.CollectionID)
	if err != nil {
		return errs.New("failed to check for existing item").WithError(err)
	}

	var e *events.Event
	if existingItem != nil {
		var existingAgg aggregate.Aggregate
		existingAgg, err = aggregate.NewWithID(aggregate.AggregateTypeItem, existingItem.ID)
		if err != nil {
			return errs.New("failed to create aggregate for existing item").WithError(err)
		}

		e, err = events.New(existingAgg, &payload.BumpItem{})
	} else {
		e, err = events.NewWithGeneratedID(
			aggregate.AggregateTypeItem,
			&payload.CreateItem{
				Content:      data,
				Mimetype:     opts.Mimetype,
				CollectionID: opts.CollectionID,
			},
		)
	}
	if err != nil {
		return errs.New("failed to create copy event: ").WithError(err)
	}

	if _, _, err = i.manager.Apply(ctx, e, events.AppendOptions{ExpectedVersion: -1}); err != nil {
		return errs.New("failed to apply copy event").WithError(err)
	}

	return nil
}

// Get retrieves the last copied data from the Bore instance.
func (b *Bore) Get(ctx context.Context, options GetClipboardOptions) (PasteResult, error) {
	if b.clipboard.Available() && options.FromSystemClipboard {
		rawContent, err := b.clipboard.Read(ctx)
		if err != nil {
			return PasteResult{}, errs.New("failed to read from system clipboard").WithError(err)
		}

		return PasteResult{Content: rawContent, Item: nil}, nil
	}

	options.CollectionID = strings.TrimSpace(options.CollectionID)
	if options.CollectionID != "" && !options.SkipCollectionCheck {
		existingCollection, err := b.repository.Collections().
			FindOne(ctx, repository.CollectionLookupOptions{
				Identifier: options.CollectionID,
				Name:       "",
			})
		if err != nil {
			return PasteResult{}, errs.New("failed to check if collection exists").WithError(err)
		} else if existingCollection == nil {
			return PasteResult{}, errs.ErrCollectionNotFound
		}
	}

	var (
		item *models.Item
		err  error
	)

	identifier := strings.TrimSpace(options.ItemID)
	if identifier == "" {
		item, err = b.repository.Items().FindLatest(ctx, options.CollectionID)
	} else {
		item, err = b.repository.Items().FindById(ctx, identifier, options.CollectionID)
	}

	if err != nil {
		return PasteResult{}, errs.New("failed to find latest item").WithError(err)
	}

	if item == nil {
		// nolint: exhaustruct
		return PasteResult{}, nil
	}

	if options.DeleteAfterPaste {
		agg, err := aggregate.NewWithID(aggregate.AggregateTypeItem, item.ID)
		if err != nil {
			return PasteResult{}, errs.New(
				"failed to create aggregate for deletion event",
			).WithError(err)
		}

		e, err := events.New(agg, &payload.DeleteItem{})
		if err != nil {
			return PasteResult{}, errs.New("failed to create delete event").WithError(err)
		}

		if _, _, err = b.manager.Apply(ctx, e); err != nil {
			return PasteResult{}, errs.New("failed to apply delete event").WithError(err)
		}
	}

	return PasteResult{Content: item.Content, Item: item}, nil
}
