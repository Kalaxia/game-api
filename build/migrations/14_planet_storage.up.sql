CREATE TABLE IF NOT EXISTS map__planet_storage(
  id SERIAL PRIMARY KEY,
  capacity int NOT NULL,
  resources json NOT NULL
);
ALTER TABLE map__planets ADD COLUMN storage_id int REFERENCES map__planet_storage(id);


CREATE OR REPLACE FUNCTION remove_planet_storage()
  RETURNS trigger AS
$BODY$
BEGIN
 DELETE FROM map__planet_storage WHERE id = OLD.storage_id;
 
 RETURN OLD;
END;
$BODY$
LANGUAGE PLPGSQL;

CREATE TRIGGER trigger_planet_storage_removal
    AFTER DELETE ON map__planets
    FOR EACH ROW
    EXECUTE PROCEDURE remove_planet_storage();