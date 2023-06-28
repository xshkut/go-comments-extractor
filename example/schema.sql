-- Generated SQL Schema

-- source: ../go/first.go

CREATE TABLE page (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    file_id INTEGER REFERENCES files(id) NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE task_page (
    task_id INTEGER NOT NULL REFERENCES tasks(id),
    page_id INTEGER NOT NULL REFERENCES pages(id),
    position INTEGER NOT NULL,
    PRIMARY KEY (task_id, page_id, position)
);

-- source: ../go/main.go

CREATE TYPE file_source AS ENUM ('aws_s3', 'google_slides', 'text');

CREATE TABLE file (
    id SERIAL PRIMARY KEY,
    source file_source NOT NULL,
    key TEXT NOT NULL,
    created_by INTEGER NOT NULL REFERENCES user(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);