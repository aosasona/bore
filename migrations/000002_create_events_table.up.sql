CREATE TABLE IF NOT EXISTS events (
	sequence_id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id TEXT NOT NULL UNIQUE,

	aggregate_type TEXT NOT NULL, -- e.g. 'item'
	aggregate_id TEXT NOT NULL,   -- e.g. valid ULID
	aggregate_version INTEGER NOT NULL, -- version of the aggregate after this event (1 to N)

	type TEXT NOT NULL,
	payload BLOB NOT NULL CHECK (json_valid(payload)),
-- use ISO8601 format for timestamps
	occurred_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

	stream_id TEXT GENERATED ALWAYS AS (aggregate_type || ':' || aggregate_id) STORED
);

-- bun:split
CREATE UNIQUE INDEX IF NOT EXISTS ux_events_stream_version ON events(aggregate_type, aggregate_id, aggregate_version);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_stream ON events(stream_id);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_sequence_id ON events(sequence_id);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_stream_seq ON events(aggregate_type, aggregate_id, sequence_id);

