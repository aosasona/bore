package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type collectionRepository struct {
	db *bun.DB
}

// Exists implements CollectionRepository.
func (c *collectionRepository) Exists(ctx context.Context, identifier string) (bool, error) {
	if identifier == "" {
		return false, nil
	}

	var flag int
	err := c.db.NewSelect().
		TableExpr("collections").
		ColumnExpr("1").
		Where("id = ?", identifier).
		Limit(1).
		Scan(ctx, &flag)
	if err != nil {
		return false, err
	}

	return flag == 1, nil
}

var _ CollectionRepository = (*collectionRepository)(nil)
