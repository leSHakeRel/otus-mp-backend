-- +goose Up
-- +goose StatementBegin
CREATE TABLE evenings (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_evenings_owner_id ON evenings(owner_id);
CREATE INDEX idx_evenings_scheduled_at ON evenings(scheduled_at);
CREATE INDEX idx_evenings_is_private ON evenings(is_private);
-- +goose StatementEnd
