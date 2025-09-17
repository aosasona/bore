DROP INDEX IF EXISTS idx_events_sequence_id;

-- bun:split
DROP INDEX IF EXISTS idx_events_type;

-- bun:split
DROP INDEX IF EXISTS idx_events_aggregate_version;

-- bun:split
DROP INDEX IF EXISTS idx_events_aggregate;

-- bun:split
DROP TABLE IF EXISTS events;

