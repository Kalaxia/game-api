ALTER TABLE players ADD faction_id int references faction__factions(id), ADD is_active boolean NOT NULL DEFAULT false;
