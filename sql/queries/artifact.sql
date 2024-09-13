-- name: UpsertArtifact :one
INSERT INTO artifacts (content, content_sha256, collection_id) VALUES (:content, sha256(:content), :collection_id)
  ON CONFLICT(content_sha256, COALESCE(collection_id, 'root'))
  DO UPDATE SET last_modified = unixepoch()
  RETURNING *;

-- name: UpdateArtifactLastModified :exec
UPDATE artifacts SET last_modified = unixepoch() WHERE id = :id;

-- name: GetLatestArtifact :one
SELECT * FROM artifacts ORDER BY last_modified DESC LIMIT 1;

-- name: DeleteAndReturnLatestArtifact :one
DELETE FROM artifacts WHERE id = (
  SELECT id FROM artifacts ORDER BY last_modified DESC LIMIT 1
) RETURNING *;

-- name: DeleteArtifactById :exec
DELETE FROM artifacts WHERE id = :id;

-- name: DeleteLatestArtifact :exec
DELETE FROM artifacts WHERE id = (
  SELECT id FROM artifacts ORDER BY last_modified DESC LIMIT 1
);
