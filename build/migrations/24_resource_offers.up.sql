CREATE TABLE IF NOT EXISTS trade__resource_offers(
    id SERIAL PRIMARY KEY,
    operation VARCHAR(15) NOT NULL,
    quantity INT NOT NULL,
    lot_quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    resource VARCHAR(25) NOT NULL,
    location_id INT NOT NULL REFERENCES map__planets(id),
    destination_id INT REFERENCES map__planets(id),
    created_at TIMESTAMPTZ NOT NULL,
    accepted_at TIMESTAMPTZ
);
