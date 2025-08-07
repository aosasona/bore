CREATE TABLE collections (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  hash TEXT NOT NULL, -- the folder names are not stored but the path hash is stored and used for lookup with "." as separator
  pinned_at TIMESTAMP, -- timestamp in seconds when the collection was pinned
  created_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
  updated_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),

  UNIQUE (name, hash)
);

