CREATE TABLE IF NOT EXISTS ship__ships(
    id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES ship__models(id) ON DELETE CASCADE,
    hangar_id INT REFERENCES map__planets(id),
    fleet_id INT REFERENCES fleet__fleets(id),
    construction_state_id INT REFERENCES map__planet_point_productions(id),
    created_at timestamptz NOT NULL
);

DROP TABLE map__planet_hangar_groups;
DROP TABLE fleet__combat_squadron_actions;
DROP TABLE fleet__combat_squadrons;
DROP TABLE fleet__combat_rounds;
DROP TABLE fleet__squadrons;
DROP TABLE ship__construction_groups;
