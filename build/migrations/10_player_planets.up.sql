ALTER TABLE map__planets ADD player_id int references players(id) ON DELETE SET NULL;
