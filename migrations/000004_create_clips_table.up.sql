CREATE TABLE IF NOT EXISTS `clips` (
  id TEXT PRIMARY KEY NOT NULL,
  content TEXT NOT NULL,
  hash TEXT NOT NULL,
  mimetype TEXT NOT NULL DEFAULT 'text/plain',
  size INTEGER NOT NULL DEFAULT 0,

  created_at TIMESTAMP NOT NULL DEFAULT (unixepoch()),
  updated_at TIMESTAMP NOT NULL DEFAULT (unixepoch()),

  collection_id TEXT,

  FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE (hash, collection_id)
);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_clips_collection_id ON clips (collection_id);

-- bun:split
CREATE INDEX IF NOT EXISTS idx_clips_hash ON clips (hash);

