CREATE TABLE IF NOT EXISTS map__planet_buildings(
  id SERIAL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  type VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  planet_id int references map__planets(id)
);
