package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/uptrace/bun"
)

type clipRepository struct {
	db *bun.DB
}

// Create implements ClipRepository.
func (c *clipRepository) Create(ctx context.Context, clip *Clip) error {
	panic("unimplemented")
}

// DeleteById implements ClipRepository.
func (c *clipRepository) DeleteById(ctx context.Context, identifier string) error {
	ctx, cancel := withContext(ctx)
	defer cancel()

	identifier = strings.TrimSpace(identifier)
	_, err := c.db.NewDelete().Model((*Clip)(nil)).Where("id = ?", identifier).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// FindById implements ClipRepository.
func (c *clipRepository) FindById(ctx context.Context, identifier string) (*Clip, error) {
	ctx, cancel := withContext(ctx)
	defer cancel()

	identifier = strings.TrimSpace(identifier)

	clip := new(Clip)
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
func (c *clipRepository) FindLatest(ctx context.Context, collectionID string) (*Clip, error) {
	ctx, cancel := withContext(ctx)
	defer cancel()

	collectionID = strings.TrimSpace(collectionID)

	clip := new(Clip)
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

var _ ClipRepository = (*clipRepository)(nil)
