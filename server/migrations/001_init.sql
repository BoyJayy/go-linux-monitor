CREATE TABLE IF NOT EXISTS devices (
    id BIGSERIAL PRIMARY KEY,
    host_id TEXT NOT NULL UNIQUE,
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS metric_snapshots (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    collected_at TIMESTAMPTZ NOT NULL,

    cpu_usage_percent DOUBLE PRECISION,
    memory_used_percent DOUBLE PRECISION,
    memory_total_bytes BIGINT,
    memory_used_bytes BIGINT,

    network_rx_bps DOUBLE PRECISION,
    network_tx_bps DOUBLE PRECISION,

    cpu JSONB NOT NULL DEFAULT '{}'::jsonb,
    memory JSONB NOT NULL DEFAULT '{}'::jsonb,
    disk JSONB NOT NULL,
    network JSONB NOT NULL,
    raw JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_metric_snapshots_device_time
ON metric_snapshots (device_id, collected_at DESC);
