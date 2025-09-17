package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/uptrace/bun"
)

type itemsRepository struct {
	db *bun.DB
}

// Create implements ClipRepository.
func (c *itemsRepository) Create(ctx context.Context, clip *Item) error {
	ctx, cancel := withContext(ctx)
	defer cancel()

	_, err := c.db.NewInsert().Model(clip).Exec(ctx)
	return err
}

// DeleteById implements ClipRepository.
func (c *itemsRepository) DeleteById(ctx context.Context, identifier string) error {
	ctx, cancel := withContext(ctx)
	defer cancel()

	identifier = strings.TrimSpace(identifier)
	_, err := c.db.NewDelete().Model((*Item)(nil)).Where("id = ?", identifier).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// FindById implements ClipRepository.
func (c *itemsRepository) FindById(ctx context.Context, identifier string) (*Item, error) {
	ctx, cancel := withContext(ctx)
	defer cancel()

	identifier = strings.TrimSpace(identifier)

	clip := new(Item)
	query := c.db.NewSelect().Model(clip).
		Where("id = ?", identifier).
		Limit(1)

	if err := query.Scan(ctx, clip); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return clip, nil
}

// GetLastClip implements ClipRepository.
func (c *itemsRepository) FindLatest(ctx context.Context, collectionID string) (*Item, error) {
	ctx, cancel := withContext(ctx)
	defer cancel()

	collectionID = strings.TrimSpace(collectionID)

	clip := new(Item)
	query := c.db.NewSelect().Model(clip).
		Where("collection_id = ?", collectionID).
		Order("created_at DESC").
		Limit(1)

	if err := query.Scan(ctx, clip); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return clip, nil
}

var _ ItemRepository = (*itemsRepository)(nil)
