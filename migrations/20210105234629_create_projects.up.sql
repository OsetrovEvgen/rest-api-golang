BEGIN;

CREATE TABLE projects (
    id varchar not null primary key,
    name varchar not null,
    description varchar not null
);

COMMIT;