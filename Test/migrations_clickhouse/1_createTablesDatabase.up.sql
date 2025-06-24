CREATE TABLE IF NOT EXISTS log_buf
(
    dt            DateTime64(9, 'Europe/Moscow'),
    api           Nullable(String),
    service_name  Nullable(String),
    server_key    Nullable(String),
    request       Nullable(String),
    response      Nullable(String),
    response_code Nullable(UInt32)
)
    engine = Buffer('default', 'log_buf', 16, 10, 100, 10000, 1000000, 10000000, 100000000);