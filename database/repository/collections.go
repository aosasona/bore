package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/pkg/errs"
	"go.trulyao.dev/bore/v2/pkg/validation"
)

type collectionRepository struct {
	db *bun.DB
}

// FindAll implements CollectionRepository.
func (c *collectionRepository) FindAll(
	ctx context.Context,
	opts FindAllOptions,
) (models.Collections, error) {
	query := c.db.NewSelect().
		Model((*models.Collection)(nil)).
		ColumnExpr("collection.*").
		ColumnExpr("(SELECT COUNT(*) FROM item WHERE item.collection_id = collection.id) AS items_count").
		Table("collection")

	for _, order := range opts.OrderBy {
		direction := "DESC"
		if order.Ascending {
			direction = "ASC"
		}
		query = query.OrderExpr(order.Field + " " + direction)
	}

	if opts.Pagination != nil {
		if opts.Pagination.Limit > 0 {
			query = query.Limit(opts.Pagination.Limit)
		}
		if opts.Pagination.Offset > 0 {
			query = query.Offset(opts.Pagination.Offset)
		}
	}

	var collections models.Collections
	if err := query.Scan(ctx, &collections); err != nil {
		return nil, err
	}

	fmt.Println(query.String()) // TODO: remove

	return collections, nil
}

// Rename implements CollectionRepository.
func (c *collectionRepository) Rename(
	ctx context.Context,
	tx bun.Tx,
	identifier string,
	newName string,
) error {
	if strings.TrimSpace(identifier) == "" {
		return ErrEmptyIdentifier
	}

	newName = strings.TrimSpace(newName)
	if !validation.IsValidCollectionName(newName) {
		return validation.ErrInvalidCollectionName
	}

	res, err := tx.NewUpdate().
		Model((*models.Collection)(nil)).
		Set("name = ?", newName).
		Where("id = ?", identifier).
		Exec(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.New("no collection found with the given identifier")
	}

	return nil
}

// DeleteById implements CollectionRepository.
func (c *collectionRepository) DeleteById(ctx context.Context, tx bun.Tx, identifier string) error {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return ErrEmptyIdentifier
	}

	_, err := tx.NewDelete().Model((*models.Collection)(nil)).Where("id = ?", identifier).Exec(ctx)
	return err
}

// FindOne implements CollectionRepository.
func (c *collectionRepository) FindOne(
	ctx context.Context,
	opts CollectionLookupOptions,
) (*models.Collection, error) {
	opts.Identifier = strings.TrimSpace(opts.Identifier)
	opts.Name = strings.TrimSpace(opts.Name)

	if opts.Identifier == "" && opts.Name == "" {
		return nil, errs.New("either identifier or name must be provided")
	}

	collection := new(models.Collection)
	query := c.db.NewSelect().Model(collection).
		Where("id = ?", opts.Identifier).
		WhereOr("LOWER(name) = LOWER(?)", opts.Name).
		Limit(1)
	err := query.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.Wrap(err, "failed to find collection")
	}

	return collection, nil
}

// FindByName implements CollectionRepository.
func (c *collectionRepository) FindByName(
	ctx context.Context,
	name string,
) (*models.Collection, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errs.New("name cannot be empty")
	}

	collection := new(models.Collection)
	err := c.db.NewSelect().
		Model(collection).
		Where("LOWER(name) = LOWER(?)", name).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return collection, nil
}

// Create implements CollectionRepository.
func (c *collectionRepository) Create(
	ctx context.Context,
	tx bun.Tx,
	collection *models.Collection,
) error {
	_, err := tx.NewInsert().Model(collection).Ignore().Exec(ctx)
	return err
}

// FindById implements CollectionRepository.
func (c *collectionRepository) FindById(
	ctx context.Context,
	identifier string,
) (*models.Collection, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, ErrEmptyIdentifier
	}

	collection := new(models.Collection)
	err := c.db.NewSelect().Model(collection).Where("id = ?", identifier).Limit(1).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return collection, nil
}

var _ CollectionRepository = (*collectionRepository)(nil)
