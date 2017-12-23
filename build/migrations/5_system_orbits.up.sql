CREATE TABLE IF NOT EXISTS map__system_orbits(
  id SERIAL PRIMARY KEY,
  radius int,
  system_id int references map__systems(id)
);
