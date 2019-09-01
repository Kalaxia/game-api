CREATE TABLE IF NOT EXISTS map__territories(
    id SERIAL PRIMARY KEY,
    map_id INT NOT NULL REFERENCES map__maps(id) ON DELETE CASCADE,
    planet_id INT NOT NULL REFERENCES map__planets(id),
    military_influence INT NOT NULL,
    political_influence INT NOT NULL,
    economic_influence INT NOT NULL,
    cultural_influence INT NOT NULL,
    religious_influence INT NOT NULL
);

CREATE TABLE IF NOT EXISTS map__territory_histories(
    id SERIAL PRIMARY KEY,
    territory_id INT NOT NULL REFERENCES map__territories(id) ON DELETE CASCADE,
    player_id INT REFERENCES players(id) ON DELETE SET NULL,
    action VARCHAR(25) NOT NULL,
    data JSON NOT NULL,
    happened_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS map__planet_territories(
    territory_id INT NOT NULL REFERENCES map__territories(id) ON DELETE CASCADE,
    planet_id INT NOT NULL REFERENCES map__planets(id) ON DELETE CASCADE,
    status VARCHAR(15) NOT NULL
);

CREATE TABLE IF NOT EXISTS map__system_territories(
    territory_id INT NOT NULL REFERENCES map__territories(id) ON DELETE CASCADE,
    system_id INT NOT NULL REFERENCES map__systems(id) ON DELETE CASCADE,
    status VARCHAR(15) NOT NULL
);