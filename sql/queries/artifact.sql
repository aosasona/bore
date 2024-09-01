-- name: GetArtifactByContent :one
SELECT * FROM artifacts WHERE content = :content
