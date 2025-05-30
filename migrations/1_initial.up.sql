CREATE TABLE wgconfigs (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    filename TEXT NOT NULL,
);