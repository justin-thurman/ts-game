-- migrate:up
CREATE INDEX idx_accounts_username ON accounts(username);

-- migrate:down
DROP INDEX idx_accounts_username;
