-- Dump the data from the artifacts table into a temporary table
CREATE TEMP TABLE temp_artifacts AS SELECT * FROM table_name;

-- Drop indexes and constraints
DROP INDEX idx_artifacts_sha256_collection ON artifacts;

-- Drop the artifacts table
DROP TABLE artifacts;

-- Recreate the artifacts table without the type column
CREATE TABLE artifacts (
  id TEXT PRIMARY KEY DEFAULT (uuid()), -- uuid
  content BLOB NOT NULL,
  content_sha256 TEXT NOT NULL,
  last_modified INTEGER NOT NULL DEFAULT (unixepoch()),

  collection_id TEXT,
  FOREIGN KEY (collection_id) REFERENCES collections(id)
) strict;

CREATE UNIQUE INDEX idx_artifacts_sha256_collection ON artifacts(content_sha256, COALESCE(collection_id, 'root'));

-- Copy the data back into the artifacts table
INSERT INTO artifacts (id, content, content_sha256, last_modified, collection_id) SELECT id, content, content_sha256, last_modified, collection_id FROM temp_artifacts;

-- Drop the temporary table
DROP TABLE temp_artifacts;
