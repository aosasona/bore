CREATE TABLE artifacts (
  id TEXT primary key default (uuid()), -- uuid
  content BLOB NOT NULL,
  content_sha256 TEXT NOT NULL,
  type TEXT NOT NULL DEFAULT 'text',
  last_modified INTEGER NOT NULL DEFAULT (unixepoch()),

  collection_id TEXT,
  FOREIGN KEY (collection_id) REFERENCES collections(id)
) strict;

CREATE UNIQUE INDEX idx_artifacts_sha256_collection ON artifacts(content_sha256, COALESCE(collection_id, 'root'));
