CREATE TABLE IF NOT EXISTS ip_allocations(
    ip           TEXT PRIMARY KEY,
    pubkey   TEXT NOT NULL,
    prefix       TEXT NOT NULL,
    allocated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    released_at  DATETIME
);

CREATE TABLE IF NOT EXISTS ip_allocator_state(
    prefix         TEXT PRIMARY KEY,
    next_candidate TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_ip_allocations_prefix_released 
    ON ip_allocations(prefix, released_at);

CREATE INDEX IF NOT EXISTS idx_ip_allocations_pubkey 
    ON ip_allocations(pubkey);