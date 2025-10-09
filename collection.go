package bore

import (
	"context"

	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events"
	"go.trulyao.dev/bore/v2/pkg/events/aggregate"
	"go.trulyao.dev/bore/v2/pkg/events/payload"
	"go.trulyao.dev/bore/v2/pkg/validation"
)

type collectionNamespace struct {
	*Bore
}

func (c *collectionNamespace) Create(ctx context.Context, name string) error {
	if !validation.IsValidCollectionName(name) {
		return validation.ErrInvalidCollectionName
	}

	existingCollection, err := c.repository.Collections().
		FindOne(ctx, repository.CollectionLookupOptions{
			Identifier: "",
			Name:       name,
		})
	if err != nil {
		return err
	}

	if existingCollection != nil {
		return errs.New("collection with the same name already exists")
	}

	event, err := events.NewWithGeneratedID(
		aggregate.AggregateTypeCollection,
		&payload.CreateCollection{Name: name},
	)
	if err != nil {
		return errs.New("failed to create collection creation event").WithError(err)
	}

	if _, _, err = c.manager.Apply(ctx, event); err != nil {
		return errs.New("failed to apply collection creation event").WithError(err)
	}

	return nil
}

// Delete deletes a collection and all its associated items.
func (c *collectionNamespace) Delete(ctx context.Context, identifier string) error {
	agg, err := aggregate.WithID(aggregate.AggregateTypeCollection, identifier)
	if err != nil {
		return errs.New("failed to create aggregate for deletion event").WithError(err)
	}

	event, err := events.New(agg, &payload.DeleteCollection{})
	if err != nil {
		return errs.New("failed to create collection deletion event").WithError(err)
	}

	if _, _, err = c.manager.Apply(ctx, event); err != nil {
		return errs.New("failed to apply collection deletion event").WithError(err)
	}

	return nil
}

// Rename renames a collection. This does not change the collection ID.
// It will return an error if the new name is already taken by another collection.
func (c *collectionNamespace) Rename(ctx context.Context, identifier, newName string) error {
	if !validation.IsValidCollectionName(newName) {
		return validation.ErrInvalidCollectionName
	}

	existingCollection, err := c.repository.Collections().FindByName(ctx, newName)
	if err != nil {
		return errs.New("failed to check existing collection name").WithError(err)
	}

	if existingCollection != nil {
		return errs.New("collection with the same name already exists")
	}

	agg, err := aggregate.WithID(aggregate.AggregateTypeCollection, identifier)
	if err != nil {
		return errs.New("failed to create aggregate for rename event").WithError(err)
	}

	event, err := events.New(agg, &payload.RenameCollection{NewName: newName})
	if err != nil {
		return errs.New("failed to create collection rename event").WithError(err)
	}

	if _, _, err = c.manager.Apply(ctx, event); err != nil {
		return errs.New("failed to apply collection rename event").WithError(err)
	}

	return nil
}

// List returns all collections with their associated items count.
func (c *collectionNamespace) List(
	ctx context.Context,
) ([]*repository.CollectionWithItemsCount, error) {
	panic("not implemented")
}
