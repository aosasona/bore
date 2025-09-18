DROP INDEX IF EXISTS idx_events_stream_seq;

-- bun:split
DROP INDEX IF EXISTS idx_events_sequence_id;

-- bun:split
DROP INDEX IF EXISTS idx_events_type;

-- bun:split
DROP INDEX IF EXISTS idx_events_stream;

-- bun:split
DROP INDEX IF EXISTS ux_events_stream_version;

-- bun:split
DROP TABLE IF EXISTS events;

