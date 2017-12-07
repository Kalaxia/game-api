CREATE TABLE IF NOT EXISTS players(
  id SERIAL PRIMARY KEY,
  username VARCHAR(180) NOT NULL UNIQUE,
  pseudo VARCHAR(180) NOT NULL UNIQUE,
  server_id int references servers(id),
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
)
