CREATE TABLE IF NOT EXISTS `clips` (
  id TEXT NOT NULL PRIMARY KEY,
  content TEXT NOT NULL,
  hash TEXT NOT NULL,
  mimetype TEXT NOT NULL DEFAULT 'text/plain',
  size INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
  updated_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
  collection_id TEXT NOT NULL,
  device_id TEXT NOT NULL,

  FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE (hash, collection_id)
);

