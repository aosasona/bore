package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type clipRepository struct {
	db *bun.DB
}

// GetLastClip implements ClipRepository.
func (c *clipRepository) GetLastClip(ctx context.Context) (Clip, error) {
	var clip Clip
	err := c.db.NewSelect().
		Model(&clip).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx, &clip)

	return clip, err
}

var _ ClipRepository = (*clipRepository)(nil)
