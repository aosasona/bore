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

// FindClipById implements ClipRepository.
func (c *clipRepository) FindClipById(ctx context.Context, identifier string) (*Clip, error) {
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
func (c *clipRepository) FindLatestClip(ctx context.Context, collectionID string) (*Clip, error) {
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
