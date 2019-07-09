ALTER TABLE faction__factions ADD wallet INT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS faction__settings(
    id SERIAL PRIMARY KEY,
    faction_id INT NOT NULL REFERENCES faction__factions(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    value INT NOT NULL,
    updated_at TIMESTAMPTZ
);