insert into markets (
  id,
  base_token_address,
  base_token_decimals,
  base_token_symbol,
  base_token_name,
  quote_token_address,
  quote_token_decimals,
  quote_token_symbol,
  quote_token_name,
  min_order_size,
  maker_fee_rate,
  taker_fee_rate,
  price_precision,
  price_decimals,
  amount_decimals,
  gas_used_estimation,
  is_published,
  created_at
) values (
  'HOT-WETH',
  '0x6829f329f8f0768ad62a65477514deEd90825564',
  18,
  'HOT',
  'HOT',
  ' 0xc778417e063141139fce010982780140aa0cd5ab',
  18,
  'WETH',
  'WETH',
  0.001,
  0.003,
  0.001,
  5,
  5,
  5,
  1,
  1,
  datetime('now')
);

insert into tokens (address, symbol, name, decimals) values
('0x6829f329f8f0768ad62a65477514deEd90825564', 'HOT', 'HOT', 18),
(' 0xc778417e063141139fce010982780140aa0cd5ab', 'WETH', 'WETH', 18);
