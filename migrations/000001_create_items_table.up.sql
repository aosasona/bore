CREATE TABLE IF NOT EXISTS items (
	id TEXT PRIMARY KEY NOT NULL,
	content TEXT NOT NULL,
	hash TEXT NOT NULL,
	mimetype TEXT NOT NULL DEFAULT 'text/plain',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_applied_sequence_id INTEGER,

-- collection is optional; if not provided, item is considered uncategorized
	collection_id TEXT REFERENCES collections(id) ON DELETE CASCADE
);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_items_collection_id ON items(collection_id);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_items_created_at ON items(created_at);

-- Ensure hash is unique within a collection (or globally if uncategorized)
CREATE UNIQUE INDEX IF NOT EXISTS idx_items_unique_collection_hash ON items(collection_id, hash);

-- Trigger to update updated_at on row modification
CREATE TRIGGER trg_update_items_updated_at
AFTER UPDATE ON items
FOR EACH ROW
BEGIN
	UPDATE items SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
end
;

