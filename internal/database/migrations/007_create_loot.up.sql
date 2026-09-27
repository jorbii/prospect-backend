CREATE TABLE loot (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    item_id UUID NOT NULL
        REFERENCES items(id)
        ON DELETE RESTRICT,

    quantity BIGINT NOT NULL DEFAULT 1,

    position_x REAL NOT NULL,
    position_y REAL NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT loot_quantity_positive
        CHECK (quantity > 0)
);