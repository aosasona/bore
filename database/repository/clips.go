package repository

import "github.com/uptrace/bun"

type clipRepository struct {
	db *bun.DB
}

// GetLastClip implements ClipRepository.
func (c *clipRepository) GetLastClip() (Clip, error) {
	panic("unimplemented")
}

var _ ClipRepository = (*clipRepository)(nil)
