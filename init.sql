PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

-- 1. Cấu hình Target
CREATE TABLE IF NOT EXISTS monitor_targets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    base_url TEXT NOT NULL,
    check_interval INTEGER DEFAULT 5,
    is_active BOOLEAN DEFAULT 1
);

-- 2. Dữ liệu nóng cho Dashboard (General View)
CREATE TABLE IF NOT EXISTS health_snapshots (
    target_id INTEGER PRIMARY KEY,
    pid INTEGER,               -- PID của chương trình đang bị soi
    last_status INTEGER,
    last_latency REAL,
    cpu_usage REAL,
    ram_usage_mb REAL,
    tcp_connections INTEGER,
    handle_count INTEGER,
    health_score INTEGER,
    bottleneck TEXT,
    error_message TEXT,        -- Lưu nhanh lỗi nếu có
    last_check_at INTEGER,     -- Unix Timestamp
    FOREIGN KEY (target_id) REFERENCES monitor_targets(id) ON DELETE CASCADE
);

-- 3. Dữ liệu vẽ biểu đồ (Lưu 7 ngày)
CREATE TABLE IF NOT EXISTS health_metrics_raw (
    target_id INTEGER NOT NULL,
    timestamp_ns INTEGER NOT NULL, 
    latency_ms REAL NOT NULL,
    status_code INTEGER,
    PRIMARY KEY (target_id, timestamp_ns)
) WITHOUT ROWID;

-- 4. Tổng kết cuối ngày
CREATE TABLE IF NOT EXISTS health_daily_summary (
    target_id INTEGER NOT NULL,
    date_unix INTEGER NOT NULL, -- Unix timestamp của 00:00:00 ngày hôm đó
    avg_latency REAL,
    uptime_percentage REAL,
    total_checks INTEGER,
    PRIMARY KEY (target_id, date_unix)
);