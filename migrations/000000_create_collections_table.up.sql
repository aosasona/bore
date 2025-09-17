CREATE TABLE IF NOT EXISTS collections (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_collections_name ON collections(name);

-- bun:split
CREATE TRIGGER IF NOT EXISTS trg_collections_updated_at
AFTER UPDATE ON collections
FOR EACH ROW
BEGIN
	UPDATE collections SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
end
;

