UPDATE exchange_lkup SET deleted_at=now() WHERE exchange_id = 1 AND deleted_at IS NULL;
INSERT INTO exchange_lkup (address, exchange_id) VALUES ('paloma1rlnkl275xqekyf5g62lda2yutvjj2rqnewr8m2q362jfvd45eyzqkptf2j', 1);
