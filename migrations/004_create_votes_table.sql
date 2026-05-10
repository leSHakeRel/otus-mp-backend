-- Migration: Create votes table
-- Version: 1.0.0

CREATE TABLE votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    evening_id UUID NOT NULL REFERENCES evenings(id) ON DELETE CASCADE,
    evening_film_id UUID NOT NULL REFERENCES evening_films(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    value INTEGER NOT NULL CHECK (value >= 1 AND value <= 5),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(evening_id, evening_film_id, user_id)
);

CREATE INDEX idx_votes_evening_id ON votes(evening_id);
CREATE INDEX idx_votes_user_id ON votes(user_id);
CREATE INDEX idx_votes_evening_film_id ON votes(evening_film_id);
