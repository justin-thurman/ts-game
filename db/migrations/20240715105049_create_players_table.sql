-- migrate:up
CREATE TABLE players (
  id SERIAL PRIMARY KEY,
  account_id INT NOT NULL,
  FOREIGN KEY (account_id) REFERENCES accounts(id),
  name VARCHAR(100) NOT NULL UNIQUE,
  class INT NOT NULL,
  stats_id INT NOT NULL UNIQUE,
  FOREIGN KEY (stats_id) REFERENCES playerstats(id),
  equipment_id INT NOT NULL UNIQUE,
  FOREIGN KEY (equipment_id) REFERENCES equipment(id),
  inventory INT[] NOT NULL,
  curr_health INT NOT NULL,
  max_health INT NOT NULL,
  curr_xp INT NOT NULL DEFAULT 0,
  player_level INT NOT NULL DEFAULT 1,
  room_id INT NOT NULL DEFAULT 1,
  recall_room_id INT NOT NULL DEFAULT 1
);

-- migrate:down
DROP TABLE players;

