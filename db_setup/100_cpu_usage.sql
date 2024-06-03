CREATE DATABASE homework;
\c homework
CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE TABLE cpu_usage(
  ts    TIMESTAMPTZ,
  host  TEXT,
  usage DOUBLE PRECISION
);
SELECT create_hypertable('cpu_usage', 'ts');
-- copies sample data into hypertable. comment out if not wanted 
COPY cpu_usage FROM '/docker-entrypoint-initdb.d/cpu_usage.csv' CSV HEADER;