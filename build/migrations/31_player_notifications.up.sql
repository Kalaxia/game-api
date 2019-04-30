CREATE TABLE IF NOT EXISTS player__notifications(
    id SERIAL PRIMARY KEY,
    player_id INT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,
    content VARCHAR(100) NOT NULL,
    data JSON NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    read_at TIMESTAMPTZ
);