CREATE TABLE IF NOT EXISTS goods_events (
    ID Int64,
    ProjectID Int64,
    Name String,
    Description String,
    Priority Int64,
    Removed Bool,
    EventTime DateTime
) ENGINE = MergeTree()
ORDER BY (EventTime, ID);