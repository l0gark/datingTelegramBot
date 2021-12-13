CREATE TABLE IF NOT EXISTS users
(
    id          varchar PRIMARY KEY NOT NULL,
    name        varchar             NOT NULL,
    sex         boolean             NOT NULL,
    age         int                 NOT NULL,
    description varchar             NOT NULL,
    city        varchar             NOT NULL,
    image       varchar             NOT NULL,
    started     bool                NOT NULL,
    stage       int                 NOT NULL,
    chat_id     int                 NOT NULL
);
