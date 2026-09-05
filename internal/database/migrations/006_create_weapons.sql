CREATE TABLE weapons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    item_id UUID NOT NULL UNIQUE
        REFERENCES items(id)
        ON DELETE CASCADE,

    damage INTEGER NOT NULL,
    fire_rate INTEGER NOT NULL,
    magazine_size INTEGER NOT NULL,
    reload_time REAL NOT NULL,
    range INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT weapons_damage_positive
        CHECK (damage > 0),

    CONSTRAINT weapons_fire_rate_positive
        CHECK (fire_rate > 0),

    CONSTRAINT weapons_magazine_size_positive
        CHECK (magazine_size > 0),

    CONSTRAINT weapons_reload_time_positive
        CHECK (reload_time > 0),

    CONSTRAINT weapons_range_positive
        CHECK (range > 0)
);