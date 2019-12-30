CREATE TABLE IF NOT EXISTS faction__wars(
    id SERIAL PRIMARY KEY,
    faction_id INT NOT NULL REFERENCES faction__factions(id) ON DELETE CASCADE,
    target_id INT NOT NULL REFERENCES faction__factions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS faction__casus_belli(
    id SERIAL PRIMARY KEY,
    faction_id INT NOT NULL REFERENCES faction__factions(id) ON DELETE CASCADE,
    victim_id INT NOT NULL REFERENCES faction__factions(id) ON DELETE CASCADE,
    player_id INT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    war_id INT DEFAULT NULL REFERENCES faction__wars(id) ON DELETE SET NULL,
    type VARCHAR(30) NOT NULL,
    data JSONB DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL
)