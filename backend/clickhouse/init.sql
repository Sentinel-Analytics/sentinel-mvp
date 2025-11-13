CREATE DATABASE IF NOT EXISTS sentinel;

CREATE TABLE IF NOT EXISTS sentinel.events (
    Timestamp DateTime,
    SiteID String,
    ClientIP String,
    URL String,
    Referrer String,
    ScreenWidth UInt16,
    Browser String,
    OS String,
    Country String
) ENGINE = MergeTree()
ORDER BY (SiteID, Timestamp);
