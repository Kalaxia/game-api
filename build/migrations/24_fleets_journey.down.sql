
ALTER TABLE fleet__journeys DROP current_step_id;
ALTER TABLE fleet__journeys_steps DROP next_step_id;

DROP TABLE fleet__journeys_steps;

ALTER TABLE fleet__fleets DROP map_pos_x;
ALTER TABLE fleet__fleets DROP map_pos_y;
