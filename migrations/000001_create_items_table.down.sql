DROP INDEX IF EXISTS idx_items_collection_id;

-- bun:split
DROP INDEX IF EXISTS idx_items_created_at;

-- bun:split
DROP INDEX IF EXISTS idx_items_unique_collection_hash;

-- bun:split
DROP TRIGGER IF EXISTS trg_update_items_updated_at;

-- bun:split
DROP TABLE IF EXISTS items;

