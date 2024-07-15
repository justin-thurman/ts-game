-- migrate:up
CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  salt VARCHAR(255) NOT NULL
);

-- migrate:down
DROP TABLE accounts;
