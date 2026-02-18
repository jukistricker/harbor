-- Lưu trữ mọi bản tin để vẽ biểu đồ và tính Delta
CREATE TABLE IF NOT EXISTS harbor_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_id TEXT,             -- Định danh (VD: main.exe_13612)
    timestamp INTEGER,          -- Unix Epoch từ Agent
    
    -- Layer 2
    handle_count INTEGER,
    thread_count INTEGER,
    mem_private INTEGER,
    total_cycles INTEGER,
    
    -- Layer 3
    active_connections INTEGER,
    
    -- Layer 4
    disk_read INTEGER,
    disk_write INTEGER,
    net_out_errors INTEGER,
    
    mode TEXT                   -- Normal / Hyper
);

-- Index để truy vấn chuỗi thời gian cực nhanh
CREATE INDEX IF NOT EXISTS idx_metrics_timeline ON harbor_metrics(target_id, timestamp);

