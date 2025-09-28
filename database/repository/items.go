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

// Bump implements ItemRepository.
func (i *itemRepository) Bump(
	ctx context.Context,
	tx bun.Tx,
	identifier string,
	sequenceId int64,
) error {
	item, err := i.FindById(ctx, identifier, "")
	if err != nil {
		return err
	}

	if item == nil {
		return errors.New("item not found")
	}

	item.LastAppliedSequenceID = sequenceId

	_, err = tx.NewUpdate().Model(item).
		Column("last_applied_sequence_id").
		WherePK().
		Exec(ctx)
	return err
}

// Create implements ItemRepository.
func (i *itemRepository) Create(ctx context.Context, tx bun.Tx, item *models.Item) error {
	_, err := tx.NewInsert().Model(item).Exec(ctx)
	return err
}

// DeleteById implements ItemRepository.
func (i *itemRepository) DeleteById(ctx context.Context, tx bun.Tx, identifier string) error {
	identifier = strings.TrimSpace(identifier)
	_, err := tx.NewDelete().Model((*models.Item)(nil)).Where("id = ?", identifier).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// FindByHash implements ItemRepository.
func (i *itemRepository) FindByHash(
	ctx context.Context,
	hash string,
	collectionId string,
) (*models.Item, error) {
	hash = strings.TrimSpace(hash)

	item := new(models.Item)
	query := i.db.NewSelect().Model(item).
		Where("hash = ?", hash)
	if collectionId != "" {
		query.Where("collection_id = ?", collectionId)
	}

	if err := query.Limit(1).Scan(ctx, item); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

// FindById implements ItemRepository.
func (i *itemRepository) FindById(
	ctx context.Context,
	identifier string,
	collectionId string,
) (*models.Item, error) {
	identifier = strings.TrimSpace(identifier)

	item := new(models.Item)
	query := i.db.NewSelect().Model(item).
		Where("id = ?", identifier)

	if collectionId != "" {
		query.Where("collection_id = ?", collectionId)
	}

	if err := query.Limit(1).Scan(ctx, item); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

// GetLastItem implements ItemRepository.
func (i *itemRepository) FindLatest(
	ctx context.Context,
	collectionID string,
) (*models.Item, error) {
	collectionID = strings.TrimSpace(collectionID)

	item := new(models.Item)
	query := i.db.NewSelect().Model(item).
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
