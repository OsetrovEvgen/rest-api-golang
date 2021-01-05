BEGIN;

CREATE TABLE comments (
    id varchar not null primary key,
    task_id varchar not null REFERENCES tasks ON DELETE CASCADE,
    text varchar not null,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;