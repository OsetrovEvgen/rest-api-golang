BEGIN;

CREATE TABLE columns (
    id varchar not null primary key,
    project_id varchar not null REFERENCES projects ON DELETE CASCADE,
    name varchar not null,
    position integer,
    UNIQUE (project_id, name)
);

COMMIT;