CREATE TABLE items (
    id bigint GENERATED ALWAYS AS IDENTITY,
    title varchar(255) NOT NULL,
    url text NOT NULL,
    note text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT item_id PRIMARY KEY (id)
)