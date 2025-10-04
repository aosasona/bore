package repository

import (
	"context"

	"github.com/uptrace/bun"
	"go.trulyao.dev/bore/v2/database/models"
	"go.trulyao.dev/bore/v2/pkg/errs"
)

type collectionRepository struct {
	db *bun.DB
}

// Create implements CollectionRepository.
func (c *collectionRepository) Create(
	ctx context.Context,
	tx bun.Tx,
	collection *models.Collection,
) error {
	_, err := tx.NewInsert().Model(collection).Exec(ctx)
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
		return nil, err
	}

	return &collection, nil
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
