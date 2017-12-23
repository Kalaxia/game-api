CREATE TABLE IF NOT EXISTS map__maps(
  id SERIAL PRIMARY KEY,
  server_id int references servers(id),
  size int
);
