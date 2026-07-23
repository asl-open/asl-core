CREATE TABLE sources (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title          TEXT NOT NULL,
    author         TEXT NOT NULL,
    type           TEXT NOT NULL,
    edition        TEXT NOT NULL,
    language       TEXT NOT NULL,
    locator_scheme TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT sources_type_check CHECK (
        type IN ('quran', 'hadith_collection', 'fiqh_manual', 'tafsir', 'article', 'other')
    )
);
