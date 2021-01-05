BEGIN;

CREATE TABLE tasks (
    id varchar not null primary key,
    column_id varchar not null REFERENCES columns ON DELETE CASCADE,
    name varchar not null,
    description varchar not null,
    position integer
);

COMMIT;