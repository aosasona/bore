CREATE TABLE IF NOT EXISTS events (
	sequence_id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id TEXT NOT NULL UNIQUE,
	aggregate TEXT, -- e.g. 'item:<item_id>'
	aggregate_version INTEGER,
	type TEXT NOT NULL,
	payload JSON NOT NULL,
-- use ISO8601 format for timestamps
	occurred_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_aggregate ON events(aggregate);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_aggregate_version ON events(aggregate, aggregate_version);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_events_sequence_id ON events(sequence_id);

