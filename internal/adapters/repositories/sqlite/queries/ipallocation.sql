-- name: GetReleasedIP :one
SELECT ip FROM ip_allocations
WHERE prefix = ? AND released_at IS NOT NULL
ORDER BY released_at ASC LIMIT 1;

-- name: ReleaseIP :exec
UPDATE ip_allocations
SET released_at = NULL, allocated_at = CURRENT_TIMESTAMP
WHERE ip = ?;

-- name: GetNextCandidate :one
SELECT next_candidate FROM ip_allocator_state WHERE prefix = ?;

-- name: UpdateNextCandidate :exec
INSERT INTO ip_allocator_state (prefix, next_candidate) VALUES (?, ?)
ON CONFLICT(prefix) DO UPDATE SET next_candidate = excluded.next_candidate;

-- name: IsAllocated :one
SELECT COUNT(*) FROM ip_allocations
WHERE ip = ? AND released_at IS NULL;

-- name: Reserve :exec
INSERT INTO ip_allocations (ip, pubkey, prefix)
VALUES (?, ?, ?)
ON CONFLICT(ip) DO UPDATE SET
    pubkey       = excluded.pubkey,
    prefix       = excluded.prefix,
    allocated_at = CURRENT_TIMESTAMP,
    released_at  = NULL;

-- name: Release :exec
UPDATE ip_allocations
SET released_at = CURRENT_TIMESTAMP
WHERE ip = ? AND released_at IS NULL;

-- name: ReleaseByPublicKey :exec
UPDATE ip_allocations 
SET released_at = CURRENT_TIMESTAMP 
WHERE pubkey = ? AND released_at IS NULL;