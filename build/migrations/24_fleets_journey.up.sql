-- add the position on the map
ALTER TABLE fleet__fleets ADD map_pos_x FLOAT DEFAULT 0;
ALTER TABLE fleet__fleets ADD map_pos_y FLOAT DEFAULT 0;

CREATE TABLE IF NOT EXISTS fleet__journeys_steps(
    id SERIAL PRIMARY KEY,
    journey_id int references fleet__journeys(id),
    --next_step_id int references fleet__journeys_steps(id), -- if unset => last step
    ---------------
    planet_start_id int references map__planets(id),
    map_pos_x_start FLOAT DEFAULT 0,
    map_pos_y_start FLOAT DEFAULT 0,
    ---------------
    planet_final_id int references map__planets(id),
    map_pos_x_final FLOAT DEFAULT 0,
    map_pos_y_final FLOAT DEFAULT 0,
    ---------------
    time_start timestamp NOT NULL,
    time_jump timestamp NOT NULL,
    time_arrival timestamp NOT NULL,
    ---------------
    step_number INT NOT NULL
);
ALTER TABLE fleet__journeys_steps ADD next_step_id int references fleet__journeys_steps(id);

ALTER TABLE fleet__journeys ADD current_step_id int references fleet__journeys_steps(id);

-- TODO remove unsed code
