
ALTER TABLE fleet__journeys DROP first_step_id;

DROP TABLE fleet__journeys_steps;

ALTER TABLE fleet__fleets DROP map_pos_x;
ALTER TABLE fleet__fleets DROP map_pos_y;
