CREATE TABLE contributors (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name TEXT NOT NULL,
    roles        TEXT[] NOT NULL,
    bio          TEXT,
    handle       TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT contributors_roles_check CHECK (
        cardinality(roles) > 0
        AND roles <@ ARRAY['contributor', 'translator', 'reviewer', 'editor']::TEXT[]
    )
);
