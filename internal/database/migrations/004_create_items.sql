CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,

    stackable BOOLEAN NOT NULL DEFAULT FALSE,
    max_stack INTEGER NOT NULL DEFAULT 1,

    rarity TEXT NOT NULL DEFAULT 'common',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT items_max_stack_positive CHECK (max_stack > 0)
);