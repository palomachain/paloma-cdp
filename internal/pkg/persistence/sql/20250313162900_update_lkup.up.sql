UPDATE exchange_lkup SET deleted_at=now() WHERE exchange_id = 1 AND deleted_at IS NULL;
INSERT INTO exchange_lkup (address, exchange_id) VALUES ('paloma13pth7njfkssewshex2cwwgqn5hyke38gtvum6qsnk9alq3dpe82qw5mtgj', 1);
