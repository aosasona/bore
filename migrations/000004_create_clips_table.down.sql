DROP INDEX IF EXISTS idx_clips_collection_id;

-- bun:split
DROP INDEX IF EXISTS idx_clips_device_id;

-- bun:split
DROP INDEX IF EXISTS idx_clips_hash;

-- bun:split
DROP TABLE IF EXISTS clips;

