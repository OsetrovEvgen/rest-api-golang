BEGIN;

CREATE TABLE projects (
    id varchar not null primary key,
    name varchar not null,
    description varchar not null
);

CREATE TABLE columns (
    id varchar not null primary key,
    project_id varchar not null REFERENCES projects ON DELETE CASCADE,
    name varchar not null,
    position integer,
    UNIQUE (project_id, name)
);

CREATE TABLE tasks (
    id varchar not null primary key,
    column_id varchar not null REFERENCES columns ON DELETE CASCADE,
    name varchar not null,
    description varchar not null,
    position integer
);

CREATE TABLE comments (
    id varchar not null primary key,
    task_id varchar not null REFERENCES tasks ON DELETE CASCADE,
    text varchar not null,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
