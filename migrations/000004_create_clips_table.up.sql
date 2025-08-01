CREATE TABLE IF NOT EXISTS `clips` (
  id TEXT PRIMARY KEY,
  content TEXT NOT NULL,
  hash TEXT NOT NULL,
  mimetype TEXT NOT NULL DEFAULT 'text/plain',
  size INTEGER NOT NULL DEFAULT 0,

  created_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
  updated_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

  collection_id TEXT,
  device_id TEXT NOT NULL,

  FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE (hash, collection_id)
);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_clips_collection_id ON clips (collection_id);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_clips_device_id ON clips (device_id);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_clips_hash ON clips (hash);

