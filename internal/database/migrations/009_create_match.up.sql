CREATE TYPE match_status AS ENUM (
    'waiting',
    'active',
    'finished'
);

CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    status match_status NOT NULL DEFAULT 'waiting',

    max_players INTEGER NOT NULL DEFAULT 10,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    started_at TIMESTAMPTZ,

    ended_at TIMESTAMPTZ,

    CONSTRAINT matches_max_players_check
        CHECK (max_players > 1)
);