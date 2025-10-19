DROP INDEX IF EXISTS uq_collection_name;

-- bun:split
CREATE UNIQUE INDEX uq_collection_name_ci ON collections (LOWER(name));

