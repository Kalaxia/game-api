DROP TABLE map__planet_production_points;

CREATE TABLE IF NOT EXISTS ship__construction_states(
    id SERIAL PRIMARY KEY,
    current_points INT NOT NULL DEFAULT 0,
    points INT NOT NULL
);
CREATE TABLE IF NOT EXISTS map__planet_construction_states(
    id SERIAL PRIMARY KEY,
    built_at timestamptz not null,
    current_points int not null,
    points int not null
);

ALTER TABLE map__planet_buildings MODIFY construction_state_id INT DEFAULT NULL REFERENCES map__planet_construction_states;
ALTER TABLE ship__ships MODIFY construction_state_id INT DEFAULT NULL REFERENCES ship__construction_states;

CREATE OR REPLACE FUNCTION remove_ship_construction_state()
  RETURNS trigger AS
$BODY$
BEGIN
 DELETE FROM ship__construction_states WHERE id = OLD.construction_state_id;
 
 RETURN OLD;
END;
$BODY$
LANGUAGE PLPGSQL;

CREATE TRIGGER trigger_ship_construction_state_removal
    AFTER DELETE ON ship__ships
    FOR EACH ROW
    EXECUTE PROCEDURE remove_ship_construction_state();

CREATE OR REPLACE FUNCTION remove_building_construction_state()
  RETURNS trigger AS
$BODY$
BEGIN
 DELETE FROM map__planet_construction_states WHERE id = OLD.construction_state_id;
 
 RETURN OLD;
END;
$BODY$
LANGUAGE PLPGSQL;

CREATE TRIGGER trigger_building_construction_state_removal
    AFTER DELETE ON map__planet_buildings
    FOR EACH ROW
    EXECUTE PROCEDURE remove_building_construction_state();