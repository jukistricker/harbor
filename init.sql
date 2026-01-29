-- CHẾ ĐỘ HIỆU NĂNG CAO
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = OFF;
PRAGMA page_size = 4096;
PRAGMA cache_size = -10000; -- 10MB RAM cache cho Index

-- 1. Bảng Cấu hình (Ít ghi, đọc nhiều)
CREATE TABLE targets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    base_url TEXT NOT NULL,
    monitor_interval INTEGER DEFAULT 5, -- Giây
    retention_days INTEGER DEFAULT 7,
    created_at INTEGER
);

CREATE TABLE target_endpoints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_id INTEGER NOT NULL,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    is_monitored BOOLEAN DEFAULT 1,
    FOREIGN KEY(target_id) REFERENCES targets(id)
);

-- 2. Bảng Metrics (CHỈ GHI - RAW DATA)
-- Sử dụng WITHOUT ROWID và PK (endpoint_id, timestamp_ns) để dữ liệu xếp chồng tuần tự trên ổ cứng
CREATE TABLE health_metrics (
    endpoint_id INTEGER NOT NULL,
    timestamp_ns INTEGER NOT NULL, -- Unix Nano
    status_code INTEGER NOT NULL,
    latency_ms REAL NOT NULL,
    
    -- Network Trace (Phân tích điểm nghẽn như Postman)
    dns_ms REAL DEFAULT 0,
    tcp_ms REAL DEFAULT 0,
    tls_ms REAL DEFAULT 0,
    ttfb_ms REAL DEFAULT 0, -- Chốt chặn Code/DB
    load_ms REAL DEFAULT 0, -- Chốt chặn Network/Size
    
    res_size_bytes INTEGER DEFAULT 0,
    test_type INTEGER DEFAULT 1, -- 1: HealthCheck, 2: Stress (Sample), 3: Single
    
    PRIMARY KEY (endpoint_id, timestamp_ns)
) WITHOUT ROWID;

-- 3. Bảng Snapshot (DỮ LIỆU NÓNG - PHỤC VỤ UI)
-- Bảng này luôn chỉ có số dòng bằng số lượng Endpoint. UI chỉ cần Query bảng này.
CREATE TABLE health_snapshots (
    endpoint_id INTEGER PRIMARY KEY,
    last_status INTEGER,
    last_latency REAL,
    
    -- Các chỉ số tính toán (Cập nhật bởi Go Worker)
    avg_latency_24h REAL DEFAULT 0,
    jitter_ms REAL DEFAULT 0,      -- Độ biến động (Tính trong RAM rồi ghi xuống)
    uptime_24h REAL DEFAULT 0,     -- Success Rate (%)
    
    -- Kết luận nhanh
    bottleneck_type TEXT,          -- 'DNS', 'Server', 'Network', 'None'
    last_check_at INTEGER
);

-- 4. Bảng Stress Test (Chỉ lưu kết quả tổng kết)
CREATE TABLE stress_reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_id INTEGER NOT NULL,
    test_at INTEGER NOT NULL,
    total_reqs INTEGER,
    concurrency INTEGER,
    p99_ms REAL,
    rps REAL,
    success_rate REAL,
    avg_ttfb_ms REAL,
    note TEXT
);