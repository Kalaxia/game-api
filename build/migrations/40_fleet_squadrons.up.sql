DROP TABLE ship__ships;

CREATE TABLE IF NOT EXISTS ship__construction_groups(
    id SERIAL PRIMARY KEY,
    location_id INT NOT NULL REFERENCES map__planets(id) ON DELETE CASCADE,
    model_id INT NOT NULL REFERENCES ship__models(id) ON DELETE CASCADE,
    construction_state_id INT NOT NULL REFERENCES map__planet_point_productions(id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS fleet__squadrons(
    id SERIAL PRIMARY KEY,
    fleet_id INT NOT NULL REFERENCES fleet__fleets(id) ON DELETE CASCADE,
    ship_model_id INT NOT NULL REFERENCES ship__models(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    combat_initiative INT DEFAULT 0,
    combat_position JSONB NOT NULL,
    position JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS fleet__combat_rounds(
    id SERIAL PRIMARY KEY,
    combat_id INT NOT NULL REFERENCES fleet__combats(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS fleet__combat_squadrons(
    id SERIAL PRIMARY KEY,
    fleet_id INT NOT NULL REFERENCES fleet__fleets(id) ON DELETE CASCADE,
    round_id INT NOT NULL REFERENCES fleet__combat_rounds(id) ON DELETE CASCADE,
    ship_model_id INT NOT NULL REFERENCES ship__models(id) ON DELETE CASCADE,
    initiative INT DEFAULT 0,
    quantity INT NOT NULL,
    position JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS fleet__combat_squadron_actions(
    id SERIAL PRIMARY KEY,
    loss INT DEFAULT 0,
    squadron_id INT NOT NULL REFERENCES fleet__combat_squadrons(id) ON DELETE CASCADE,
    target_id INT DEFAULT NULL REFERENCES fleet__combat_squadrons(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS map__planet_hangar_groups(
    id SERIAL PRIMARY KEY,
    location_id INT NOT NULL REFERENCES map__planets(id) ON DELETE CASCADE,
    model_id INT NOT NULL REFERENCES ship__models(id) ON DELETE CASCADE,
    quantity INT NOT NULL
);