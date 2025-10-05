package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/pkg/errs"
)

type collectionRepository struct {
	db *bun.DB
}

type CollectionWithItemsCount struct {
	models.Collection
	ItemsCount int64 `bun:"items_count"`
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

	var collection *models.Collection
	err := c.db.NewSelect().Model(&collection).
		WherePK(opts.Identifier).
		WhereOr("LOWER(name) = LOWER(?)", opts.Name).
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

// FindByName implements CollectionRepository.
func (c *collectionRepository) FindByName(
	ctx context.Context,
	name string,
) (*models.Collection, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errs.New("name cannot be empty")
	}

	var collection models.Collection
	err := c.db.NewSelect().
		Model(&collection).
		Where("LOWER(name) = LOWER(?)", name).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &collection, nil
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
	if identifier == "" {
		return nil, errs.New("identifier cannot be empty")
	}

	var collection models.Collection
	err := c.db.NewSelect().Model(&collection).WherePK(identifier).Limit(1).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &collection, nil
}

var _ CollectionRepository = (*collectionRepository)(nil)
