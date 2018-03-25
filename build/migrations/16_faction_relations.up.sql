CREATE TABLE IF NOT EXISTS diplomacy__factions(
  faction_id int references faction__factions(id),
  other_faction_id int references faction__factions(id),
  state VARCHAR(15) NOT NULL
);
