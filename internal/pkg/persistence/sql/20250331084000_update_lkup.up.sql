UPDATE exchange_lkup SET deleted_at=now() WHERE exchange_id = 2 AND deleted_at IS NULL;
INSERT INTO exchange_lkup (address, exchange_id) VALUES ('paloma198dta03m4epzanh9kf9wda9e3g7lcll7v6tmnjdnf7g8y49h7k4sjg8l7l', 2);
