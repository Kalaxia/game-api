CREATE TABLE IF NOT EXISTS players(
  id SERIAL PRIMARY KEY,
  username VARCHAR(180) NOT NULL,
  pseudo VARCHAR(180) NOT NULL,
  server_id int NOT NULL references servers(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  CONSTRAINT username_constraint UNIQUE (username, server_id),
  CONSTRAINT pseudo_constraint UNIQUE (pseudo, server_id)
);
