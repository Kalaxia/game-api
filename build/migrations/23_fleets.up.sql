CREATE TABLE IF NOT EXISTS fleet__journeys(
  id SERIAL PRIMARY KEY,
  created_at timestamp,
  ended_at timestamp
);

CREATE TABLE IF NOT EXISTS fleet__fleets(
  id SERIAL PRIMARY KEY,
  player_id int NOT NULL references players(id) ON DELETE CASCADE,
  location_id int references map__planets(id),
  journey_id int references fleet__journeys(id)
);
ALTER TABLE ship__ships ADD fleet_id int references fleet__fleets(id);


CREATE OR REPLACE FUNCTION remove_fleet_journey()
  RETURNS trigger AS
$BODY$
BEGIN
 DELETE FROM fleet_journeys WHERE id = OLD.journey_id;
 
 RETURN OLD;
END;
$BODY$
LANGUAGE PLPGSQL;

CREATE TRIGGER trigger_fleet_removal
    AFTER DELETE ON fleet__fleets
    FOR EACH ROW
    EXECUTE PROCEDURE remove_fleet_journey();