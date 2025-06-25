CREATE TABLE IF NOT EXISTS goods_events (
    Id Int64,
    ProjectId Int64,
    Name String,
    Description String,
    Priority Int64,
    Removed Bool,
    EventTime DateTime
) ENGINE = MergeTree()
ORDER BY (EventTime, Id);