-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    evening_id UUID NOT NULL REFERENCES evenings(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_evening_id ON comments(evening_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
-- +goose StatementEnd
