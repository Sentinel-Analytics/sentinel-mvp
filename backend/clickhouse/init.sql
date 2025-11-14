CREATE DATABASE IF NOT EXISTS sentinel;

CREATE TABLE IF NOT EXISTS events (
    Timestamp DateTime,
    SiteID String,
    ClientIP String,
    URL String,
    Referrer String,
    ScreenWidth UInt16,
    Browser String,
    OS String,
    Country String,
    TrustScore UInt8,
    LCP Nullable(Float64),
    CLS Nullable(Float64),
    FID Nullable(Float64)
) ENGINE = MergeTree()
ORDER BY (SiteID, Timestamp);

CREATE TABLE IF NOT EXISTS session_events (
    Timestamp DateTime,
    SiteID String,
    SessionID String,
    Payload String
) ENGINE = MergeTree()
ORDER BY (SiteID, SessionID, Timestamp);
