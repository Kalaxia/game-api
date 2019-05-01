CREATE TABLE IF NOT EXISTS fleet__combats(
  id SERIAL PRIMARY KEY,
  attacker_id int NOT NULL references fleet__fleets(id) ON DELETE CASCADE,
  defender_id int NOT NULL references fleet__fleets(id) ON DELETE CASCADE,
  attacker_ships json not null,
  defender_ships json not null,
  attacker_losses json,
  defender_losses json,
  is_victory boolean not null,
  begin_at timestamptz not null,
  end_at timestamptz
);