-- migrate:up
CREATE TABLE equipment (
  id SERIAL PRIMARY KEY,
  body INT NOT NULL DEFAULT 0,
  legs INT NOT NULL DEFAULT 0,
  helm INT NOT NULL DEFAULT 0,
  main_weapon INT NOT NULL DEFAULT 0
);

-- migrate:down
DROP TABLE equipment;
