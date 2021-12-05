CREATE TABLE IF NOT EXISTS likes
(
    id      varchar PRIMARY KEY NOT NULL,
    from_id varchar             NOT NULL REFERENCES users (id),
    to_id   varchar             NOT NULL REFERENCES users (id),
    showed  bool                NOT NULL
);
