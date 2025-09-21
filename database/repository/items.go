package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
)

type itemRepository struct {
	db *bun.DB
}

// Create implements ItemRepository.
func (c *itemRepository) Create(ctx context.Context, item *models.Item) error {
	_, err := c.db.NewInsert().Model(item).Exec(ctx)
	return err
}

// DeleteById implements ItemRepository.
func (c *itemRepository) DeleteById(ctx context.Context, identifier string) error {
	identifier = strings.TrimSpace(identifier)
	_, err := c.db.NewDelete().Model((*models.Item)(nil)).Where("id = ?", identifier).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// FindById implements ItemRepository.
func (c *itemRepository) FindById(ctx context.Context, identifier string) (*models.Item, error) {
	identifier = strings.TrimSpace(identifier)

	item := new(models.Item)
	query := c.db.NewSelect().Model(item).
		Where("id = ?", identifier).
		Limit(1)

	if err := query.Scan(ctx, item); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

// GetLastItem implements ItemRepository.
func (c *itemRepository) FindLatest(
	ctx context.Context,
	collectionID string,
) (*models.Item, error) {
	collectionID = strings.TrimSpace(collectionID)

	item := new(models.Item)
	query := c.db.NewSelect().Model(item).
		Where("collection_id = ?", collectionID).
		Order("created_at DESC").
		Limit(1)

	if err := query.Scan(ctx, item); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

var _ ItemRepository = (*itemRepository)(nil)
