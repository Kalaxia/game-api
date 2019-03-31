CREATE TABLE IF NOT EXISTS map__planets(
  id SERIAL PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  type VARCHAR(15) NOT NULL,
  system_id int references map__systems(id) ON DELETE CASCADE,
  orbit_id int references map__system_orbits(id)
);
