CREATE TABLE IF NOT EXISTS diplomacy__relations(
  planet_id int references map__planets(id) ON DELETE CASCADE,
  faction_id int references faction__factions(id) ON DELETE CASCADE,
  player_id int references players(id) ON DELETE CASCADE,
  score int NOT NULL DEFAULT 0
);
CREATE INDEX planet_faction_relation ON diplomacy__relations (planet_id, faction_id);
CREATE INDEX planet_player_relation ON diplomacy__relations (planet_id, player_id);
