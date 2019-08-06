CREATE TABLE IF NOT EXISTS map__planet_building_compartments(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    building_id INT NOT NULL REFERENCES map__planet_buildings(id) ON DELETE CASCADE,
    construction_state_id INT NOT NULL REFERENCES map__planet_point_productions(id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);