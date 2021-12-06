CREATE TABLE IF NOT EXISTS users
(
    id          varchar PRIMARY KEY NOT NULL,
    name        varchar             NOT NULL,
    sex         boolean             NOT NULL,
    age         int                 NOT NULL,
    description text                NOT NULL,
    city        text                NOT NULL,
    image       text                NOT NULL
);
