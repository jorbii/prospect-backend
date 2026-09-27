CREATE TABLE loadouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    player_id UUID NOT NULL UNIQUE
        REFERENCES players(id)
        ON DELETE CASCADE,

    primary_weapon_id UUID
        REFERENCES weapons(id)
        ON DELETE SET NULL,

    secondary_weapon_id UUID
        REFERENCES weapons(id)
        ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);