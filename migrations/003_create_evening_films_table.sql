-- +goose Up
-- +goose StatementBegin
CREATE TABLE evening_films (
    id UUID PRIMARY KEY,
    evening_id UUID NOT NULL REFERENCES evenings(id) ON DELETE CASCADE,
    tmdb_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    poster_path VARCHAR(255),
    backdrop_path VARCHAR(255),
    release_date DATE,
    vote_average DECIMAL(3,1),
    overview TEXT,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(evening_id, tmdb_id)
);

CREATE INDEX idx_evening_films_evening_id ON evening_films(evening_id);
CREATE INDEX idx_evening_films_tmdb_id ON evening_films(tmdb_id);
-- +goose StatementEnd
