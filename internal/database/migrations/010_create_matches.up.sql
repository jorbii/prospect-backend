CREATE TABLE matches_players (
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,

    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (match_id, player_id)
);