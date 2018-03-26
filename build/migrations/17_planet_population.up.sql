CREATE TABLE IF NOT EXISTS map__planet_settings(
    id SERIAL PRIMARY KEY,
    services_points int NOT NULL,
    building_points int NOT NULL,
    military_points int NOT NULL,
    research_points int NOT NULL
);
ALTER TABLE map__planets
    ADD COLUMN population INT NOT NULL DEFAULT 0,
    ADD COLUMN settings_id int REFERENCES map__planet_settings(id);
