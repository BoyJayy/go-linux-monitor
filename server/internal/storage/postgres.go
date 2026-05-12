package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"monitoring/api"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) SaveMetrics(ctx context.Context, m api.Metrics) error {
	var deviceID int64

	err := s.pool.QueryRow(ctx, `
		INSERT INTO devices (host_id)
		VALUES ($1)
		ON CONFLICT (host_id)
		DO UPDATE SET last_seen_at = now()
		RETURNING id
	`, m.HostID).Scan(&deviceID)
	if err != nil {
		return fmt.Errorf("upsert device: %w", err)
	}

	cpuJSON, err := json.Marshal(m.CPU)
	if err != nil {
		return fmt.Errorf("marshal cpu: %w", err)
	}

	memoryJSON, err := json.Marshal(m.Mem)
	if err != nil {
		return fmt.Errorf("marshal memory: %w", err)
	}

	diskJSON, err := json.Marshal(m.Disk)
	if err != nil {
		return fmt.Errorf("marshal disk: %w", err)
	}

	networkJSON, err := json.Marshal(m.Network)
	if err != nil {
		return fmt.Errorf("marshal network: %w", err)
	}

	rawJSON, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal raw metrics: %w", err)
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO metric_snapshots (
			device_id,
			collected_at,
			cpu_usage_percent,
			memory_used_percent,
			memory_total_bytes,
			memory_used_bytes,
			network_rx_bps,
			network_tx_bps,
			cpu,
			memory,
			disk,
			network,
			raw
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9::jsonb, $10::jsonb, $11::jsonb, $12::jsonb, $13::jsonb
		)
	`,
		deviceID,
		m.Timestamp,
		m.CPU.TotalPct,
		m.Mem.UsedPct,
		m.Mem.TotalBytes,
		m.Mem.UsedBytes,
		m.Network.RxBpsTotal,
		m.Network.TxBpsTotal,
		string(cpuJSON),
		string(memoryJSON),
		string(diskJSON),
		string(networkJSON),
		string(rawJSON),
	)
	if err != nil {
		return fmt.Errorf("insert metric snapshot: %w", err)
	}

	return nil
}

func (s *Storage) GetLastMetrics(ctx context.Context) (api.Metrics, error) {
	var rawJSON []byte

	err := s.pool.QueryRow(ctx, `
		SELECT raw
		FROM metric_snapshots
		ORDER BY created_at DESC
		LIMIT 1
	`).Scan(&rawJSON)
	if err != nil {
		return api.Metrics{}, fmt.Errorf("get last metrics: %w", err)
	}

	var metrics api.Metrics
	if err := json.Unmarshal(rawJSON, &metrics); err != nil {
		return api.Metrics{}, fmt.Errorf("unmarshal last metrics: %w", err)
	}

	return metrics, nil
}

func (s *Storage) GetListDevices(ctx context.Context) ([]Device, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, host_id, first_seen_at, last_seen_at
		FROM devices
		ORDER BY last_seen_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	defer rows.Close()
	now := time.Now()
	threshold := 30 * time.Second
	devices := make([]Device, 0)
	for rows.Next() {
		var d Device
		if err := rows.Scan(
			&d.ID,
			&d.HostID,
			&d.FirstSeenAt,
			&d.LastSeenAt,
		); err != nil {
			return nil, fmt.Errorf("scan devices: %w", err)
		}
		d.Online = now.Sub(d.LastSeenAt) <= threshold
		devices = append(devices, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate devices: %w", err)
	}
	return devices, nil
}

func (s *Storage) GetSnapByDevice(ctx context.Context) (api.Metrics, error) {
	var metric api.Metrics
	err := s.pool.QueryRow(ctx, `
		SELECT *
		from metric_snapshot
		
	`).Scan(&metric)
	if err != nil {
		return api.Metrics{}, fmt.Errorf("scan metric: %w", err)
	}

}
