ALTER TABLE fleet__journeys DROP start_place_id, DROP end_place_id;

ALTER TABLE fleet__journeys_steps
    ADD planet_start_id INT REFERENCES map__planets(id),
    ADD map_pos_x_start FLOAT DEFAULT 0,
    ADD map_pos_y_start FLOAT DEFAULT 0,
    ADD planet_final_id INT REFERENCES map__planets(id),
    ADD map_pos_x_final FLOAT DEFAULT 0,
    ADD map_pos_y_final FLOAT DEFAULT 0,
    DROP start_place_id,
    DROP end_place_id;

ALTER TABLE fleet__fleets
    ADD location_id INT NOT NULL REFERENCES map__planets(id),
    ADD map_pos_x FLOAT DEFAULT 0,
    ADD map_pos_y FLOAT DEFAULT 0,
    DROP place_id;

ALTER TABLE fleet__combats DROP place_id;

DROP TABLE map__places;