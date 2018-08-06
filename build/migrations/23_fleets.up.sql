CREATE TABLE IF NOT EXISTS fleet__journeys(
  id SERIAL PRIMARY KEY,
  created_at timestamp,
  ended_at timestamp
);

CREATE TABLE IF NOT EXISTS fleet__fleets(
  id SERIAL PRIMARY KEY,
  player_id int NOT NULL references players(id),
  location_id int references map__planets(id),
  journey_id int references fleet__journeys(id)
);
ALTER TABLE ship__ships ADD fleet_id int references fleet__fleets(id);
