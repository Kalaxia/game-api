CREATE TABLE IF NOT EXISTS map__planet_settings(
    id SERIAL PRIMARY KEY,
    services_points int NOT NULL,
    building_points int NOT NULL,
    military_points int NOT NULL,
    research_points int NOT NULL
);
ALTER TABLE map__planets
    ADD COLUMN population INT NOT NULL DEFAULT 0,
    ADD COLUMN settings_id int REFERENCES map__planet_settings(id);


CREATE OR REPLACE FUNCTION remove_planet_settings()
  RETURNS trigger AS
$BODY$
BEGIN
 DELETE FROM map__planet_settings WHERE id = OLD.settings_id;
 
 RETURN OLD;
END;
$BODY$
LANGUAGE PLPGSQL;

CREATE TRIGGER trigger_planet_settings_removal
    AFTER DELETE ON map__planets
    FOR EACH ROW
    EXECUTE PROCEDURE remove_planet_settings();