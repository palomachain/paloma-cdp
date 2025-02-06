SELECT 
  exchange_id, 
  time_bucket('1 second', time) AS bucket, 
  min(price) AS low, 
  max(price) as high, 
  first(price,time) as open, 
  last(price,time) as close
FROM price_data
WHERE 
  symbol_id=16 
  AND time > '2025-01-23' 
  AND time <= '2025-01-25'
GROUP BY exchange_id,bucket
ORDER BY exchange_id ASC, bucket ASC;
