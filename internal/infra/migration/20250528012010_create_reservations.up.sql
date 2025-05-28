CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE reservations (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	desk_id UUID NOT NULL REFERENCES desks(id),
	user_id UUID NOT NULL REFERENCES users(id),
	start_time TIMESTAMP NOT NULL,
	end_time TIMESTAMP NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending', 'confirmed', 'cancelled')) DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
