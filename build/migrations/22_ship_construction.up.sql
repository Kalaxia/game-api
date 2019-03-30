CREATE TABLE IF NOT EXISTS ship__construction_states(
    id SERIAL PRIMARY KEY,
    current_points INT NOT NULL DEFAULT 0,
    points INT NOT NULL
);

CREATE TABLE IF NOT EXISTS ship__ships(
    id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES ship__models(id) ON DELETE CASCADE,
    hangar_id INT REFERENCES map__planets(id),
    construction_state_id INT REFERENCES ship__construction_states(id),
    created_at timestamptz NOT NULL
);

DELETE FROM ship__player_models;
DELETE FROM ship__slots;
DELETE FROM ship__models;
ALTER TABLE ship__models ADD price JSON NOT NULL;

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