-- migrate:up
CREATE TABLE playerstats (
  id SERIAL PRIMARY KEY,
  str INT NOT NULL DEFAULT 0,
  con INT NOT NULL DEFAULT 0
);

-- migrate:down
DROP TABLE playerstats;
