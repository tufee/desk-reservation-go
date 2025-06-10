CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE desks (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	number INTEGER NOT NULL
);

-- Seed data
INSERT INTO desks (number) 
VALUES (1), (2), (3), (4), (5), (6), (7), (8), (9), (10);