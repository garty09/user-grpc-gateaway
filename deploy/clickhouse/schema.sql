CREATE TABLE IF NOT EXISTS users (
  timestamp UInt64,
  userId Nullable(UInt64)
) ENGINE = Kafka SETTINGS
            kafka_broker_list = 'kafka:9092',
            kafka_topic_list = 'users',
            kafka_group_name = 'statistics',
            kafka_format = 'JSONEachRow',
            kafka_num_consumers = 2;


CREATE TABLE IF NOT EXISTS users_stats (
  timestamp UInt64,
  userId Nullable(UInt64)
) ENGINE = MergeTree()
ORDER BY timestamp;

CREATE MATERIALIZED VIEW IF NOT EXISTS users_consumer TO users_stats
  AS SELECT * FROM users;
