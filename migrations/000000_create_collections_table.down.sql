DROP INDEX IF EXISTS idx_collections_name;

-- bun:split
DROP TRIGGER IF EXISTS trg_collections_updated_at;

-- bun:split
DROP TABLE IF EXISTS collections;

