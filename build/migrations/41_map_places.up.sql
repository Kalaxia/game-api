CREATE TABLE IF NOT EXISTS map__places(
    id SERIAL PRIMARY KEY,
    planet_id INT DEFAULT NULL REFERENCES map__planets(id),
    coordinates JSONB 
);

ALTER TABLE fleet__journeys
    ADD start_place_id INT NOT NULL REFERENCES map__places(id),
    ADD end_place_id INT NOT NULL REFERENCES map__places(id);

ALTER TABLE fleet__journeys_steps
    DROP planet_start_id,
    DROP map_pos_x_start,
    DROP map_pos_y_start,
    DROP planet_final_id,
    DROP map_pos_x_final,
    DROP map_pos_y_final,
    ADD start_place_id INT NOT NULL REFERENCES map__places(id),
    ADD end_place_id INT NOT NULL REFERENCES map__places(id);

ALTER TABLE fleet__fleets
    DROP location_id,
    DROP map_pos_x,
    DROP map_pos_y,
    ADD place_id INT DEFAULT NULL REFERENCES map__places(id);

ALTER TABLE fleet__combats ADD place_id INT NOT NULL REFERENCES map__places(id);