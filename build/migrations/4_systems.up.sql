CREATE TABLE IF NOT EXISTS map__systems(
  id SERIAL PRIMARY KEY,
  map_id int references map__maps(id) ON DELETE CASCADE,
  x int,
  y int
);
