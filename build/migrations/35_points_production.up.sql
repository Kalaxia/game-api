CREATE TABLE IF NOT EXISTS map__planet_point_productions(
    id SERIAL PRIMARY KEY,
    current_points INT NOT NULL,
    points INT NOT NULL
);

UPDATE map__planet_buildings SET status = 'operational' WHERE construction_state_id IS NOT NULL;

ALTER TABLE map__planet_buildings DROP construction_state_id, ADD construction_state_id INT DEFAULT NULL REFERENCES map__planet_point_productions;
ALTER TABLE ship__ships DROP construction_state_id, ADD construction_state_id INT DEFAULT NULL REFERENCES map__planet_point_productions;

DROP TRIGGER trigger_ship_construction_state_removal ON ship__ships;
DROP TRIGGER trigger_building_construction_state_removal ON map__planet_buildings;
DROP FUNCTION remove_ship_construction_state();
DROP FUNCTION remove_building_construction_state();
DROP TABLE map__planet_construction_states;
DROP TABLE ship__construction_states;