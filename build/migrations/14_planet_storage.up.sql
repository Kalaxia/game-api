CREATE TABLE IF NOT EXISTS map__planet_storage(
  id SERIAL PRIMARY KEY,
  capacity int NOT NULL,
  resources json NOT NULL
);
ALTER TABLE map__planets ADD COLUMN storage_id int REFERENCES map__planet_storage(id);
