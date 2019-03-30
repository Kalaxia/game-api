CREATE TABLE IF NOT EXISTS diplomacy__factions(
  faction_id int references faction__factions(id) ON DELETE CASCADE,
  other_faction_id int references faction__factions(id) ON DELETE CASCADE,
  state VARCHAR(15) NOT NULL
);
