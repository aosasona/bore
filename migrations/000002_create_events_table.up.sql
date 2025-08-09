CREATE TABLE IF NOT EXISTS events (
	id         TEXT PRIMARY KEY NOT NULL,
	aggregate_id TEXT, -- the id of the aggregate this event belongs to
	action     TEXT NOT NULL,
	version    INTEGER NOT NULL DEFAULT 1,
	payload    BLOB NOT NULL,
	ingested_at TIMESTAMP NOT NULL DEFAULT (unixepoch()),
	UNIQUE (aggregate_id, version) -- ensure that each event for an aggregate has a unique version
);

-- bun:split
CREATE INDEX IF NOT EXISTS events_agg_idx ON events (action, ingested_at);

-- bun:split
CREATE INDEX IF NOT EXISTS events_ts_idx ON events (ingested_at);

