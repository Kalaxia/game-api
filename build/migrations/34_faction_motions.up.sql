CREATE TABLE IF NOT EXISTS faction__motions(
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    is_approved BOOLEAN NOT NULL DEFAULT false,
    is_processed BOOLEAN NOT NULL DEFAULT false,
    data JSON NOT NULL,
    faction_id INT NOT NULL REFERENCES faction__factions(id) ON DELETE CASCADE,
    author_id INT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS faction__votes(
    id SERIAL PRIMARY KEY,
    motion_id INT NOT NULL REFERENCES faction__motions(id) ON DELETE CASCADE,
    author_id INT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    option INT NOT NULL,
    created_at TIMESTAMPTZ
);