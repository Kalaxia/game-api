CREATE TABLE IF NOT EXISTS map__planet_construction_states(
    id SERIAL PRIMARY KEY,
    built_at timestamptz not null,
    current_points int not null,
    points int not null
);
ALTER TABLE map__planet_buildings
    DROP COLUMN built_at,
    ADD COLUMN construction_state_id int REFERENCES map__planet_construction_states(id);

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