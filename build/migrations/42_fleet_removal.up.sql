ALTER TABLE fleet__fleets ADD deleted_at TIMESTAMPTZ DEFAULT NULL, ADD created_at TIMESTAMPTZ NOT NULL;
