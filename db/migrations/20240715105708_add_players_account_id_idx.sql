-- migrate:up
CREATE INDEX account_id_idx ON players(account_id);

-- migrate:down
DROP INDEX account_id_idx;
