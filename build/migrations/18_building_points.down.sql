ALTER TABLE map__planet_buildings
    ADD COLUMN built_at timestamptz,
    DROP COLUMN construction_state_id;
DROP TABLE map__planet_construction_states;
