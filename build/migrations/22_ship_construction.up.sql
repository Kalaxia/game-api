CREATE TABLE IF NOT EXISTS ship__construction_states(
    id SERIAL PRIMARY KEY,
    current_points INT NOT NULL DEFAULT 0,
    points INT NOT NULL
);

CREATE TABLE IF NOT EXISTS ship__ships(
    id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES ship__models(id),
    hangar_id INT REFERENCES map__planets(id),
    construction_state_id INT NOT NULL REFERENCES ship__construction_states(id),
    created_at timestamptz NOT NULL
);

DELETE FROM ship__player_models;
DELETE FROM ship__slots;
DELETE FROM ship__models;
ALTER TABLE ship__models ADD price JSON NOT NULL;
