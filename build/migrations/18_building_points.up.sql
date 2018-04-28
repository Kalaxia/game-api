CREATE TABLE IF NOT EXISTS map__planet_construction_states(
    id SERIAL PRIMARY KEY,
    built_at timestamptz not null,
    current_points int not null,
    points int not null
);
ALTER TABLE map__planet_buildings
    DROP COLUMN built_at,
    ADD COLUMN construction_state_id int REFERENCES map__planet_construction_states(id);
