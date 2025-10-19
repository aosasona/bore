package payload

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/database/repository"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/events/action"
)

type CreateCollection struct {
	Name string `json:"name" validate:"required,collection_name"`
}

// ApplyProjection implements Payload.
func (c *CreateCollection) ApplyProjection(
	ctx context.Context,
	tx bun.Tx,
	repo repository.Repository,
	options ProjectionOptions,
) error {
	if !options.Aggregate.IsValid() {
		return errs.New("invalid aggregate")
	}

	existingCollectionById, err := repo.Collections().FindById(ctx, options.Aggregate.ID())
	if err != nil {
		return err
	}

	if existingCollectionById != nil {
		// Collection already exists, nothing to do.
		return nil
	}

	/*
		NOTE: we do this to avoid name collisions

		Say Device A creates "My Collection"
		Then Device B creates "My Collection"
		Then Device B syncs to Device A

		We want to rename Device B's collection to "My Collection 001" and still keep the ID intact
		This is to avoid confusion for the user

		This way, we can keep the collections and apply any other events that have happened to the correct collection running into trouble with trying to replay events for a collectin that doesn't exists because we either failed to create it or replaced it with another collection
	*/
	existingCollection, err := repo.Collections().FindByName(ctx, c.Name)
	if err != nil {
		return err
	}

	// If a collection with the same name exists, append a number to the name.
	name := c.Name
	i := 1
	maxIterations := 1000 // Prevent infinite loops
	for existingCollection != nil {
		if i > maxIterations {
			return errs.New(
				"too many collections with the same name exist, cannot create a new one",
			)
		}

		name = fmt.Sprintf("%s %03d", c.Name, i)
		existingCollection, err = repo.Collections().FindByName(ctx, name)
		if err != nil {
			return err
		}
		i++
	}

	row := models.Collection{
		ID:   options.Aggregate.ID(),
		Name: name,
	}

	return repo.Collections().Create(ctx, tx, &row)
}

// Type implements Payload.
func (c *CreateCollection) Type() action.Action {
	return action.ActionCreateCollection
}

var _ Payload = (*CreateCollection)(nil)
