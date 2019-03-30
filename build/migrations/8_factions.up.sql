CREATE TABLE IF NOT EXISTS faction__factions(
  id SERIAL PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  description TEXT NOT NULL,
  server_id int references servers(id) ON DELETE CASCADE
);
