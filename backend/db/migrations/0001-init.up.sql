-- tokens table
create table tokens(
 symbol text primary key,
 name text not null,
 address text not null,
 decimals integer not null,
 updated_at timestamp,
 created_at timestamp
);
create unique index idx_tokens_address on tokens (address);

-- markets table
create table markets(
 id text primary key,
 base_token_address text not null,
 base_token_decimals text not null,
 base_token_symbol text not null,
 base_token_name text not null,

 quote_token_address text not null,
 quote_token_decimals  text not null,
 quote_token_symbol text not null,
 quote_token_name text not null,

 min_order_size numeric(32,18) not null,
 maker_fee_rate numeric(10,5) not null,
 taker_fee_rate numeric(10,5) not null,
 price_precision integer not null,
 price_decimals integer not null,
 amount_decimals integer not null,
 gas_used_estimation integer not null,
 is_published boolean not null default true,
 updated_at timestamp,
 created_at timestamp
);

-- trades table
create table trades(
  id SERIAL PRIMARY KEY,
  transaction_id integer not null,
  transaction_hash text,
  status text not null,
  market_id text not null,
  maker  text not null,
  taker  text not null,
  price numeric(32,18) not null,
  amount numeric(32,18) not null,
  taker_side text not null,
  maker_order_id  text not null,
  taker_order_id text not null,
  sequence int not null default 0,
  executed_at timestamp,
  updated_at timestamp,
  created_at timestamp
);
create index idx_trades_transaction_hash on trades (transaction_hash);
create index idx_trades_taker on trades (taker,market_id);
create index idx_trades_maker on trades (maker,market_id);
create index idx_market_id_status_executed_at on trades (market_id, status, executed_at);

-- orders table
create table orders(
  id text not null primary key,
  trader_address text not null,
  market_id text not null,
  side text not null,
  price  numeric(32,18) not null,
  amount  numeric(32,18) not null,
  status text not null,
  type text not null,
  version text not null,
  available_amount  numeric(32,18) not null,
  confirmed_amount  numeric(32,18) not null,
  canceled_amount  numeric(32,18) not null,
  pending_amount  numeric(32,18) not null,
  maker_fee_rate  numeric(10,5) not null,
  taker_fee_rate  numeric(10,5) not null,
  maker_rebate_rate  numeric(10,5) not null,
  gas_fee_amount  numeric(32,18) not null,
  json text not null,
  updated_at  timestamp,
  created_at  timestamp
);
create index idx_market_id_status on orders (market_id, status);
create index idx_market_trader_address on orders (trader_address, market_id, status, created_at);

-- transactions table
create table transactions(
  id SERIAL PRIMARY KEY,
  transaction_hash text,
  market_id text not null,
  status text not null,
  executed_at timestamp,
  updated_at  timestamp,
  created_at timestamp
);
create unique index idx_transactions_transaction_hash on transactions (transaction_hash);

-- launch_logs table
create table launch_logs(
  id SERIAL PRIMARY KEY,
  item_type text not null,
  item_id integer not null,
  status text not null,
  transaction_hash text,
  block_number integer,
  t_from text not null,
  t_to text not null,
  value numeric(32,18),
  gas_limit integer,
  gas_used integer,
  gas_price numeric(32,18),
  nonce integer,
  data text not null,
  executed_at timestamp,
  updated_at  timestamp,
  created_at  timestamp
);
create index idx_launch_logs_nonce on launch_logs (nonce);
create index idx_created_at on launch_logs (created_at);
create unique index idx_launch_logs_transaction_hash on launch_logs (transaction_hash);