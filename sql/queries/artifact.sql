-- name: GetArtifactByContent :one
SELECT * FROM artifacts WHERE content_sha256 = sha256(:content) AND collection_id = :collection_id;

-- name: CreateArtifact :one
INSERT INTO artifacts (content, content_sha256, type, collection_id) VALUES (:content, sha256(:content), :type, :collection_id) RETURNING *;

-- name: UpdateArtifactLastModified :exec
UPDATE artifacts SET last_modified = unixepoch() WHERE id = :id;
