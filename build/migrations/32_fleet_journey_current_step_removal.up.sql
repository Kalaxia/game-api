ALTER TABLE fleet__journeys
    DROP CONSTRAINT fleet__journeys_current_step_id_fkey,
    ADD CONSTRAINT fleet__journeys_current_step_id_fkey
        FOREIGN KEY (current_step_id)
        REFERENCES fleet__journeys_steps(id)
        ON DELETE SET NULL;