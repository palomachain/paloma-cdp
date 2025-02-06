CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS exchanges (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

--bun:split

CREATE TABLE IF NOT EXISTS symbols (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    short_name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    description TEXT NOT NULL,
    exchange_id INT NOT NULL,
    FOREIGN KEY (exchange_id) REFERENCES exchanges(id)
);

--bun:split

CREATE TABLE IF NOT EXISTS price_data (
    symbol_id INT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    FOREIGN KEY (symbol_id) REFERENCES symbols(id),
    UNIQUE (symbol_id, time)
);

--bun:split

CREATE INDEX price_data_timestamp ON price_data(
  symbol_id, 
  time
);

--bun:split
SELECT create_hypertable('price_data', by_range('time'));

--bun:split

ALTER TABLE price_data SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'symbol_id'
);

--bun:split

SELECT add_compression_policy('price_data', INTERVAL '7 days');

