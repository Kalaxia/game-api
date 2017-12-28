CREATE TABLE IF NOT EXISTS map__planet_resources(
  name VARCHAR(60) NOT NULL,
  density INT NOT NULL,
  planet_id int references map__planets(id),
  PRIMARY KEY(planet_id, name)
);
